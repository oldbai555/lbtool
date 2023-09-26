package etcd

import (
	"context"
	"errors"
	"fmt"
	"github.com/oldbai555/lbtool/pkg/dispatch"
	"github.com/oldbai555/lbtool/pkg/lock"

	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/pkg/routine"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"sync"
	"time"
)

const KeyPrefix = "baix_"

var _ dispatch.IDispatch = (*Dispatch)(nil)

type Service struct {
	*dispatch.Service
	version int64
}

type Dispatch struct {
	*lock.MulElemMuFactory

	isWatched     bool
	timeout       time.Duration
	srvMap        map[string]*Service
	unwatchSignCh chan struct{}
	etcd          *clientv3.Client
	onSrvUpdate   dispatch.OnSrvUpdatedFunc
	mu            sync.RWMutex
}

func NewDispatch(timeout time.Duration, etcdCfg clientv3.Config) (dispatch.IDispatch, error) {
	d := &Dispatch{
		timeout:       timeout,
		srvMap:        make(map[string]*Service),
		unwatchSignCh: make(chan struct{}),
	}
	d.MulElemMuFactory = lock.NewMulElemMuFactory()
	var err error
	d.etcd, err = clientv3.New(etcdCfg)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return d, nil
}

func (d *Dispatch) tryAddTimeoutToCtx(ctx context.Context) (triedCtx context.Context, cancel context.CancelFunc, hasCancel bool) {
	if d.timeout > 0 {
		triedCtx, cancel = context.WithTimeout(ctx, d.timeout)
		hasCancel = true
		return
	}
	triedCtx = ctx
	return
}

func (d *Dispatch) atomicPersistSrv(ctx context.Context, srvName string, version int64, srv *dispatch.Service) (bool, error) {
	srvJson, err := json.Marshal(srv)
	if err != nil {
		return false, err
	}

	key := d.genSvrEtcdKey(srvName)
	tx := d.etcd.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), "=", version))
	if len(srv.Nodes) == 0 {
		tx.Then(clientv3.OpDelete(key))
	} else {
		tx.Then(clientv3.OpPut(key, string(srvJson)))
	}
	resp, err := tx.Commit()
	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}

func (d *Dispatch) genSvrEtcdKey(svr string) string {
	return fmt.Sprintf("%s%s", KeyPrefix, svr)
}

func (d *Dispatch) Watch() {
	if d.isWatched {
		return
	}
	d.mu.Lock()
	if d.isWatched {
		d.mu.Unlock()
		return
	}
	d.isWatched = true
	d.mu.Unlock()

	evtCh := d.etcd.Watch(context.Background(), KeyPrefix, clientv3.WithPrefix())
	for {
		select {
		case <-d.unwatchSignCh:
			goto out
		case resp := <-evtCh:
			for _, evt := range resp.Events {
				key := string(evt.Kv.Key)
				if !strings.HasPrefix(key, KeyPrefix) {
					continue
				}

				srvName := key[len(KeyPrefix):]

				d.mu.Lock()
				oneSrvMu := d.MakeOrGetSpecElemMu(srvName)
				d.mu.Unlock()

				oneSrvMu.Lock()

				srv, ok := d.srvMap[srvName]
				if ok {
					if srv.version > evt.Kv.ModRevision {
						oneSrvMu.Unlock()
						continue
					}
				}

				switch evt.Type {
				case clientv3.EventTypePut:
					cfg := &dispatch.Service{}
					err := json.Unmarshal(evt.Kv.Value, cfg)
					if err != nil {
						delete(d.srvMap, srvName)
					} else {
						d.srvMap[srvName] = &Service{
							Service: cfg,
							version: evt.Kv.ModRevision,
						}
					}
					oneSrvMu.Unlock()
					if d.onSrvUpdate != nil {
						d.onSrvUpdate(context.Background(), dispatch.EvtUpdated, cfg)
					}
				case clientv3.EventTypeDelete:
					delete(d.srvMap, srvName)
					oneSrvMu.Unlock()
					if d.onSrvUpdate != nil {
						d.onSrvUpdate(context.Background(), dispatch.EvtDeleted, &dispatch.Service{
							SrvName: srvName,
						})
					}
				}
			}
		}
	}

out:
	if !d.isWatched {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.isWatched = false
	d.srvMap = map[string]*Service{}
}

func (d *Dispatch) LoadAll(ctx context.Context) ([]*dispatch.Service, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.isWatched {
		routine.GoV2(func() error {
			d.Watch()
			return nil
		})
	}

	newCtx, cancel, hasCancel := d.tryAddTimeoutToCtx(ctx)
	if hasCancel {
		defer cancel()
	}

	resp, err := d.etcd.Get(newCtx, KeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var services []*dispatch.Service
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		if !strings.HasPrefix(key, KeyPrefix) {
			continue
		}

		srvName := key[len(KeyPrefix):]

		srv := &dispatch.Service{}
		err = json.Unmarshal(kv.Value, srv)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal service(%s)'s value failed, err:%s", srvName, err)
		}

		d.srvMap[srvName] = &Service{
			Service: srv,
			version: kv.ModRevision,
		}

		d.onSrvUpdate(ctx, dispatch.EvtUpdated, srv)

		services = append(services, srv)
	}

	return services, nil
}

func (d *Dispatch) Register(ctx context.Context, srvName string, node *dispatch.Node) error {
	if srvName == "" {
		return errors.New("invalid serviceName, empty")
	}

	if node.Host == "" || node.Port == 0 {
		return errors.New(fmt.Sprintf("invalid node, node %+v", node))
	}

	ctx, cancel, hasCancel := d.tryAddTimeoutToCtx(ctx)
	if hasCancel {
		defer cancel()
	}

	key := d.genSvrEtcdKey(srvName)
	retry := 0
	maxRetry := 3
	for ; retry < maxRetry; retry++ {
		resp, err := d.etcd.Get(ctx, key)
		if err != nil {
			return err
		}

		srv := &dispatch.Service{
			SrvName: srvName,
		}
		var srvVersion int64
		if len(resp.Kvs) > 0 {
			val := resp.Kvs[0].Value
			err = json.Unmarshal(val, srv)
			if err != nil {
				return err
			}
			srvVersion = resp.Kvs[0].Version
		}

		// check node existed
		existed := false
		for _, n := range srv.Nodes {
			if n.Host != node.Host || n.Port != node.Port {
				continue
			}

			if !n.Available() {
				n.Status = dispatch.NodeStateAlive
			}
			n.Extra = node.Extra

			existed = true
			break
		}
		if !existed {
			srv.Nodes = append(srv.Nodes, node)
		}

		// 原子性更新
		ok, err := d.atomicPersistSrv(ctx, srvName, srvVersion, srv)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		break
	}

	if retry == maxRetry {
		return errors.New(fmt.Sprintf("set conflicted and retry fail, key %s", key))
	}
	return nil
}

func (d *Dispatch) UnRegister(ctx context.Context, srvName string, node *dispatch.Node, remove bool) error {
	if srvName == "" {
		return errors.New("invalid serviceName, empty")
	}

	ctx, cancel, hasCancel := d.tryAddTimeoutToCtx(ctx)
	if hasCancel {
		defer cancel()
	}

	key := d.genSvrEtcdKey(srvName)
	retry := 0
	maxRetry := 3
	for ; retry < maxRetry; retry++ {
		resp, err := d.etcd.Get(ctx, key)
		if err != nil {
			return err
		}

		srv := &dispatch.Service{}
		if len(resp.Kvs) == 0 {
			break
		}

		err = json.Unmarshal(resp.Kvs[0].Value, srv)
		if err != nil {
			return err
		}

		existed := false
		var remainedNodes []*dispatch.Node
		for _, n := range srv.Nodes {
			if n.Host != node.Host || n.Port != node.Port {
				remainedNodes = append(remainedNodes, n)
				continue
			}

			existed = true
			if remove {
				continue
			}
			n.Status = dispatch.NodeStateDead
			n.Extra = node.Extra
			remainedNodes = append(remainedNodes, n)
		}

		if !existed {
			return nil
		}

		srv.Nodes = remainedNodes

		ok, err := d.atomicPersistSrv(ctx, srvName, resp.Kvs[0].Version, srv)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		break
	}

	if retry == maxRetry {
		return errors.New(fmt.Sprintf("set conflicted and retry fail, key %s", key))
	}

	return nil
}

func (d *Dispatch) UnRegisterAll(ctx context.Context, srvName string) error {
	if srvName == "" {
		return errors.New("invalid serviceName, empty")
	}

	ctx, cancel, hasCancel := d.tryAddTimeoutToCtx(ctx)
	if hasCancel {
		defer cancel()
	}

	_, err := d.etcd.Delete(ctx, d.genSvrEtcdKey(srvName))
	if err != nil {
		return err
	}

	return nil
}

func (d *Dispatch) Discover(ctx context.Context, srvName string) (*dispatch.Service, error) {
	d.mu.RLock()
	srv, ok := d.srvMap[srvName]
	d.mu.RUnlock()
	if ok {
		return srv.Service, nil
	}

	d.mu.Lock()
	oneSrvMu := d.MakeOrGetSpecElemMu(srvName)
	d.mu.Unlock()

	oneSrvMu.Lock()
	defer oneSrvMu.Unlock()

	srv, ok = d.srvMap[srvName]
	if ok {
		return srv.Service, nil
	}

	if !d.isWatched {
		routine.GoV2(func() error {
			d.Watch()
			return nil
		})
	}

	ctx, cancel, hasCancel := d.tryAddTimeoutToCtx(ctx)
	if hasCancel {
		defer cancel()
	}

	resp, err := d.etcd.Get(ctx, d.genSvrEtcdKey(srvName))
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, dispatch.ErrSrvNotFound
	}

	srvCfg := &dispatch.Service{}
	err = json.Unmarshal(resp.Kvs[0].Value, srvCfg)
	if err != nil {
		return nil, err
	}

	d.srvMap[srvName] = &Service{
		Service: srvCfg,
		version: resp.Kvs[0].ModRevision,
	}

	return srvCfg, nil
}

func (d *Dispatch) OnSrvUpdated(updatedFunc dispatch.OnSrvUpdatedFunc) {
	d.onSrvUpdate = updatedFunc
}

func (d *Dispatch) UnWatch() {
	if d.unwatchSignCh == nil {
		return
	}
	d.unwatchSignCh <- struct{}{}
}
