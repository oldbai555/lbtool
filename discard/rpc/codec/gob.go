package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	enc  *gob.Encoder // 编码
	dec  *gob.Decoder // 解码
}

// 检测 GobCodec 是否实现了 Codec 接口
// 1） _ 为了避免变量未使用编译的时候报错
// 2）_ 的类型为 Codec ，接口的值为 GobCodec 的地址，(nil)表示该地址为nil。
// 这种定义方式主要用于在源码编译的时候。
var _ Codec = (*GobCodec)(nil)

// Close 关闭连接流
func (g *GobCodec) Close() error {
	return g.conn.Close()
}

// ReadHeader 读请求头
func (g *GobCodec) ReadHeader(h *Header) error {
	return g.dec.Decode(h)
}

// ReadBody 读请求体
func (g *GobCodec) ReadBody(body interface{}) error {
	return g.dec.Decode(body)
}

// Write 响应结果
func (g *GobCodec) Write(h *Header, body interface{}) error {
	defer func() {
		err := g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()

	if err := g.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}

	if err := g.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}

	return nil
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		enc:  gob.NewEncoder(buf),
		dec:  gob.NewDecoder(conn),
	}
}
