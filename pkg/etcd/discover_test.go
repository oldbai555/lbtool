package etcd

import (
	"github.com/oldbai555/lbtool/log"
	"testing"
	"time"
)

func TestNewServiceDiscovery(t *testing.T) {
	config := GetConfig()
	ser := NewServiceDiscovery(config.GetEndpointList())
	defer func(ser *ServiceDiscovery) {
		err := ser.Close()
		if err != nil {
			return
		}
	}(ser)
	err := ser.WatchService("/web/")
	if err != nil {
		return
	}
	err = ser.WatchService("/gRPC/")
	if err != nil {
		return
	}
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Infof("service list is %v", ser.GetServices())
		}
	}
}
