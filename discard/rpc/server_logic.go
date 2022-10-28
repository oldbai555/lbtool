package rpc

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"net"
	"net/http"
)

// DefaultServer is the default instance of *Server.
var DefaultServer = NewServer()

// ServerRun accepts connections on the listener and serves requests
// for each incoming connection.
func ServerRun(port uint32) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	//DefaultServer.Accept(lis)
	err = http.Serve(lis, nil)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	return nil
}

// Register publishes the receiver's methods in the DefaultServer.
func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}
