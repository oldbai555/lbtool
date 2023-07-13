package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/oldbai555/lbtool/log"
	"time"
)

// ServiceRegister 创建租约注册服务
type ServiceRegister struct {
	cli           *clientv3.Client                        //etcd client
	leaseID       clientv3.LeaseID                        //租约ID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse //租约 keepALive 相应chan
	key           string                                  //key
	val           string                                  //value
}

// NewServiceRegister 新建注册服务
func NewServiceRegister(endpoints []string, key, val string, lease int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	ser := &ServiceRegister{
		cli: cli,
		key: key,
		val: val,
	}

	//申请租约设置时间keepalive
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}

	return ser, nil
}

// 设置租约
func (s *ServiceRegister) putKeyWithLease(lease int64) error {

	//设置租约时间
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	//注册服务并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	//设置续租 定期发送需求请求
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)

	if err != nil {
		return err
	}

	s.leaseID = resp.ID
	log.Infof("lease id is %v", s.leaseID)

	s.keepAliveChan = leaseRespChan
	log.Infof("Put key:%s  val:%s  success!", s.key, s.val)
	return nil
}

// ListenLeaseRespChan 监听 续租情况
func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		log.Infof("续约成功 %v", leaseKeepResp)
		return
	}
	log.Infof("关闭续租")
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}

	log.Infof("撤销租约")
	return s.cli.Close()
}
