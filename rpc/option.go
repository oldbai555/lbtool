package rpc

import (
	"github.com/oldbai555/lb/rpc/codec"
	"time"
)

const MagicNumber = 0x3bef5c
const DefaultTimeOut = time.Second * 10

// Option 设置固定长度的字节以及编码方法来和客户端进行交互
type Option struct {
	MagicNumber    int           // MagicNumber marks this's a rpc request
	CodecType      codec.Type    // client may choose different Codec to encode body
	ConnectTimeout time.Duration // 0 means no limit
	HandleTimeout  time.Duration
}

var DefaultOption = &Option{
	MagicNumber:    MagicNumber,
	CodecType:      codec.GobType,
	ConnectTimeout: DefaultTimeOut,
	HandleTimeout:  DefaultTimeOut,
}
