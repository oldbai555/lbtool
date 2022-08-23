package codec

import "io"

type Header struct {
	ServiceMethod string `json:"service_method"` // format "Service.Method"
	Seq           uint64 `json:"seq"`            // sequence number chosen by client
	Error         string `json:"error"`
}

type Codec interface {
	io.Closer
	ReadHeader(h *Header) error
	ReadBody(body interface{}) error
	Write(h *Header, body interface{}) error
}
