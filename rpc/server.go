package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oldbai555/lb/rpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

const DefaultBufferSize = 10

// invalidRequest is a placeholder for response argv when error occurs
var invalidRequest = struct{}{}

// Server represents an RPC Server.
type Server struct {
	serviceMap sync.Map
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{}
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.
func (s *Server) Accept(lis net.Listener) {
	// for 循环等待 socket 连接建立，并开启子协程处理，处理过程交给了 ServerConn 方法。
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go s.ServeConn(conn)
	}
}

// ServeConn runs the server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() {
		_ = conn.Close()
	}()

	// 协议协商
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}

	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}

	// 得到协商的 连接的消息体协议
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}

	// 拿到消息编码对象
	cc := f(conn)

	// 对连接来的数据进行编解码处理
	s.serveCodec(cc, &opt)
}

// serveCodec 对连接来的数据进行编解码处理 进行三步逻辑
// 读取请求 readRequest
// 处理请求 handleRequest
// 回复请求 sendResponse
func (s *Server) serveCodec(cc codec.Codec, opt *Option) {
	// make sure to send a complete response
	// 回复请求的报文必须是逐个发送的，并发容易导致多个回复报文交织在一起，客户端无法解析。在这里使用锁(sending)保证。
	sending := new(sync.Mutex)

	// wait until all request are handled
	wg := new(sync.WaitGroup)

	for {
		// 读取请求
		req, err := s.readRequest(cc)
		if err != nil {
			if req == nil {
				break // it's not possible to recover, so close the connection
			}
			req.h.Error = err.Error()
			// 回复请求
			s.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}

		wg.Add(1)

		// todo timeout
		// 处理请求
		go s.handleRequest(cc, req, sending, wg, opt.HandleTimeout)
	}

	wg.Wait()
	_ = cc.Close()
}

func (s *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

// readRequest 读取请求
func (s *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := s.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}

	req.svc, req.mtype, err = s.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}

	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()

	// make sure that argvi is a pointer, ReadBody need a pointer as parameter
	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}

	if err = cc.ReadBody(argvi); err != nil {
		log.Println("rpc server: read body err:", err)
		return req, err
	}

	return req, nil
}

// sendResponse回复请求
func (s *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

// handleRequest 回复请求 - 协程并发执行请求
// 需要确保 sendResponse 仅调用一次，因此将整个过程拆分为 called 和 sent 两个阶段，在这段代码中只会发生如下两种情况：
// 1. called 信道接收到消息，代表处理没有超时，继续执行 sendResponse。
// 2. time.After() 先于 called 接收到消息，说明处理已经超时，called 和 sent 都将被阻塞。
//    在 case <-time.After(timeout) 处调用 sendResponse。
func (s *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration) {

	defer func() {
		wg.Done()
	}()

	called := make(chan struct{}, DefaultBufferSize)

	sent := make(chan struct{}, DefaultBufferSize)

	go func() {
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
			return
		}

		s.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		sent <- struct{}{}
	}()

	if timeout == 0 {
		<-called
		<-sent
		return
	}

	select {
	case <-time.After(timeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", timeout)
		s.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called:
		<-sent
	}

}

// Register publishes in the server the set of methods of the
func (s *Server) Register(rcvr interface{}) error {
	ser := newService(rcvr)
	if _, dup := s.serviceMap.LoadOrStore(ser.name, ser); dup {
		return errors.New("rpc: service already defined: " + ser.name)
	}
	return nil
}

func (s *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	// serviceMethod = Service.Method
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
		return
	}

	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svci, ok := s.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}

	svc = svci.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}

	return
}
