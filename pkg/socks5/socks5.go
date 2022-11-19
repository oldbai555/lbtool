package socks5

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// socks5 文档 https://datatracker.ietf.org/doc/html/rfc1928

const (
	Socks5Version = 0x05
	ReservedField = 0x00
)

var (
	ErrVersionNotSupported         = fmt.Errorf("protocol version not supported")
	ErrCommandNotSupported         = fmt.Errorf("request commend not supported")
	ErrReservedFieldNotSupported   = fmt.Errorf("request reserved field not supported")
	ErrAddressTypeNotSupported     = fmt.Errorf("request address type not supported")
	ErrPasswordVersionNotSupported = fmt.Errorf("password version not supported")
	ErrPasswordVerificationFailed  = fmt.Errorf("password verification failed")
	ErrMethodNotSupport            = fmt.Errorf("method not supported")
)

type Server interface {
	Run() error
}

type Socks5Server struct {
	Ip   string  `json:"ip"`
	Port int     `json:"port"`
	Conf *Config `json:"conf"`
}

type Config struct {
	AuthMethod      Method
	PasswordChecker func(username, password string) bool
	TcpTimeout      time.Duration
}

func (s *Socks5Server) Run() error {
	// localhost:1080
	address := fmt.Sprintf("%s:%d", s.Ip, s.Port)

	// 监听 tcp 端口
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("listen tcp address %s fail,err is %v", address, err)
		return err
	}

	for {
		// 接收连接
		conn, aErr := listen.Accept()
		if aErr != nil {
			log.Printf("connect failure form %s: err is %v", conn.RemoteAddr(), err)
			continue
		}

		// 处理连接
		go func() {
			// 释放连接
			defer func() {
				cErr := conn.Close()
				if cErr != nil {
					log.Printf("close connect failure form %s: err is %v", conn.RemoteAddr(), err)
				}
			}()
			hErr := s.HandleConnection(conn)
			if hErr != nil {
				log.Printf("handle connection failure form %s: err is %v", conn.RemoteAddr(), err)
			}
		}()
	}
}

// HandleConnection 处理连接
func (s *Socks5Server) HandleConnection(conn net.Conn) error {
	// 协商过程
	if err := auth(conn, s.Conf); err != nil {
		return err
	}
	// 请求过程
	if err := s.Request(conn); err != nil {
		return err
	}
	return nil
}

// Request 请求过程
func (s *Socks5Server) Request(conn io.ReadWriter) error {
	message, err := NewClientRequestMessage(conn)
	if err != nil {
		return err
	}

	if message.Atyp == AddressTypeIPv6 {
		WriteRequestFailureMessage(conn, ReplyTypeAddressTypeNotSupported)
		return ErrAddressTypeNotSupported
	}

	// check if the command is support
	switch message.Cmd {
	case CmdConnect:
		s.HandleTcp(conn, message)
	case CmdUDP:
		s.HandleUdp(conn, message)
	default:
		WriteRequestFailureMessage(conn, ReplyTypeCommandNotSupported)
		return ErrCommandNotSupported
	}

	return nil
}

// HandleTcp 处理tcp转发
func (s *Socks5Server) HandleTcp(conn io.ReadWriter, message *ClientRequestMessage) error {
	// 请求访问目标TCP服务
	address := fmt.Sprintf("%s:%d", message.Address, message.Port)
	targetConn, err := net.DialTimeout("tcp", address, s.Conf.TcpTimeout)
	if err != nil {
		WriteRequestFailureMessage(conn, ReplyTypeConnectionRefused)
		return err
	}

	// 通知成功
	addr := targetConn.LocalAddr().(*net.TCPAddr)
	if err = WriteRequestSuccessMessage(conn, addr.IP, uint16(addr.Port)); err != nil {
		return err
	}

	// 转发过程
	return TcpForward(conn, targetConn)
}

// HandleUdp 处理udp转发
func (s *Socks5Server) HandleUdp(conn io.ReadWriter, message *ClientRequestMessage) error {
	return nil
}
