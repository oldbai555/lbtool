package dispatch

import (
	"context"
	"errors"
	"math/rand"
)

var (
	ErrSrvNotFound  = errors.New("service not found")
	ErrNodeNotFound = errors.New("node not found")
)

const (
	NodeStateNil   = 0
	NodeStateAlive = 1
	NodeStateDead  = 2
)

func NewNode(host string, port int) *Node {
	return &Node{
		Host: host,
		Port: port,
	}
}

type Node struct {
	Host      string `json:"ip"`
	Port      int    `json:"port"`
	Status    int    `json:"status"`
	Priority  int    `json:"priority"`
	Name      string `json:"name"`
	SlaveFlag int    `json:"slave_flag"`
	Extra     string `json:"extra"`
}

func (n *Node) Available() bool {
	return (n.Status & NodeStateDead) == 0
}

type Service struct {
	SrvName string  `json:"srv_name"`
	Nodes   []*Node `json:"nodes"`
}

type Evt int

const (
	EvtNil Evt = iota
	EvtDeleted
	EvtUpdated
)

type OnSrvUpdatedFunc func(ctx context.Context, evt Evt, srv *Service)

type IDispatch interface {
	LoadAll(ctx context.Context) ([]*Service, error)
	Register(ctx context.Context, srvName string, node *Node) error
	UnRegister(ctx context.Context, srvName string, node *Node, remove bool) error
	UnRegisterAll(ctx context.Context, srvName string) error
	Discover(ctx context.Context, srvName string) (*Service, error)
	OnSrvUpdated(OnSrvUpdatedFunc)
	Watch()
	UnWatch()
}

func Route(ctx context.Context, discover IDispatch, srvName string) (*Node, error) {
	srv, err := discover.Discover(ctx, srvName)
	if err != nil {
		return nil, err
	}

	if len(srv.Nodes) == 0 {
		return nil, ErrNodeNotFound
	}

	nodeNum := len(srv.Nodes)
	randN := rand.Int()
	for i := 0; i < nodeNum; i++ {
		node := srv.Nodes[randN%nodeNum]
		if node.Available() {
			return node, nil
		}
		randN++
	}

	return nil, ErrNodeNotFound
}
