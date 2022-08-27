package rpc

// Call represents an active RPC.
// 承载一次 RPC 调用所需要的信息
type Call struct {
	Seq           string
	ServiceMethod string      // format "<service>.<method>"
	Args          interface{} // arguments to the function
	Reply         interface{} // reply from the function
	Error         error       // if error occurs, it will be set
	Done          chan *Call  // Strobes when call is complete.
}

func (call *Call) done() {
	call.Done <- call
}
