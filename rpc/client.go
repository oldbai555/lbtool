package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/rpc/codec"
	"github.com/oldbai555/lbtool/utils"
	"io"
	"sync"
)

var _ io.Closer = (*Client)(nil)

var ErrShutdown = errors.New("connection is shut down")

// Client represents an RPC Client.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
type Client struct {
	cc      codec.Codec
	opt     *Option
	sending sync.Mutex // 保证请求的有序发送

	header codec.Header
	mu     sync.Mutex // 保证方法有序执行

	pending map[string]*Call // 存储未处理完的请求 key call.seq

	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}

// Close the connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
	}()

	if c.closing {
		return ErrShutdown
	}

	c.closing = true
	return c.cc.Close()
}

// IsAvailable return true if the client does work
func (c *Client) IsAvailable() bool {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
	}()

	return !c.shutdown && !c.closing
}

// registerCall 将参数 call 添加到 client.pending 中
func (c *Client) registerCall(call *Call) (string, error) {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
	}()

	if c.closing || c.shutdown {
		return "", ErrShutdown
	}

	call.Seq = utils.GetRandomString(12, utils.RandomStringModNumberPlusLetter)
	c.pending[call.Seq] = call
	return call.Seq, nil
}

// removeCall 根据 seq，从 client.pending 中移除对应的 call，并返回
func (c *Client) removeCall(seq string) *Call {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
	}()

	call := c.pending[seq]

	delete(c.pending, seq)

	return call
}

// terminateCalls 服务端或客户端发生错误时调用，将 shutdown 设置为 true，且将错误信息通知所有 pending 状态的 call
func (c *Client) terminateCalls(err error) {
	c.sending.Lock()
	defer func() {
		c.sending.Unlock()
	}()

	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
	}()

	c.shutdown = true
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
}

// receive 接收结果
func (c *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = c.cc.ReadHeader(&h); err != nil {
			break
		}
		call := c.removeCall(h.Seq)

		switch {
		case call == nil:
			// call 不存在，可能是请求没有发送完整，或者因为其他原因被取消，但是服务端仍旧处理了
			err = c.cc.ReadBody(nil)

		case h.Error != "":
			//call 存在，但服务端处理出错，即 h.Error 不为空
			call.Error = fmt.Errorf(h.Error)
			err = c.cc.ReadBody(nil)
			call.done()

		default:
			//call 存在，服务端处理正常，那么需要从 body 中读取 Reply 的值
			err = c.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()

		}
	}

	// error occurs, so terminateCalls pending calls
	c.terminateCalls(err)
}

// send 发送请求
func (c *Client) send(call *Call) {
	// make sure that the c will send a complete request
	c.sending.Lock()
	defer func() {
		c.sending.Unlock()
	}()

	// register this call.
	seq, err := c.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	// prepare request header
	c.header.ServiceMethod = call.ServiceMethod
	c.header.Seq = seq
	c.header.Error = ""

	// encode and send the request
	if wErr := c.cc.Write(&c.header, call.Args); wErr != nil {
		log.Errorf("err is %v", wErr)
		rCall := c.removeCall(seq)

		// call may be nil, it usually means that Write partially failed,
		// c has received the response and handled
		if rCall != nil {
			rCall.Error = err
			rCall.done()
		}
	}

}

// Go  RPC 服务调用接口-异步接口
// It returns the Call structure representing the invocation.
func (c *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		panic(any("rpc c: done channel is unbuffered"))
	}
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	c.send(call)
	return call
}

// Call RPC 服务调用接口-同步接口
func (c *Client) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	call := c.Go(serviceMethod, args, reply, make(chan *Call, 1))
	select {
	case <-ctx.Done():
		c.removeCall(call.Seq)
		return errors.New("rpc client: call failed: " + ctx.Err().Error())
	case cCall := <-call.Done:
		return cCall.Error
	}
}
