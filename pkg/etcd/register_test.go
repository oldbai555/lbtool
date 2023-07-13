package etcd

import (
	"github.com/oldbai555/lbtool/log"
	"testing"
)

func TestNewServiceRegister(t *testing.T) {
	config := GetConfig()
	ser, err := NewServiceRegister(config.GetEndpointList(), "/web/node1", "localhost:8000", 5)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	//监听续租相应chan
	go ser.ListenLeaseRespChan()
	select {
	// case <-time.After(20 * time.Second):
	// 	ser.Close()
	}
}
