package codec

import "io"

// Header rpc 头部
type Header struct {
	ServiceMethod string `json:"service_method"` // format "Service.Method"
	Seq           uint64 `json:"seq"`            // sequence number chosen by client
	Error         string `json:"error"`
}

// Codec 消息编码抽象接口
type Codec interface {
	io.Closer
	ReadHeader(h *Header) error
	ReadBody(body interface{}) error
	Write(h *Header, body interface{}) error
}
