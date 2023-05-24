package socks5

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Method = byte

const (
	MethodNoAuth       Method = 0x00
	MethodGssApi       Method = 0x01
	MethodPassWord     Method = 0x02
	MethodNoAcceptable Method = 0xff
)

const (
	PasswordMethodVersion = 0x01
	PasswordAuthSuccess   = 0x00
	PasswordAuthFailure   = 0x01
)

type ClientAuthMessage struct {
	Version  byte     `json:"version"`
	NMethods byte     `json:"n_methods"`
	Methods  []Method `json:"methods"`
}

type ClientPasswordMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewClientAuthMessage 读取数据报文生成对应结构体
func NewClientAuthMessage(conn io.Reader) (*ClientAuthMessage, error) {

	// read version , nMethods
	buf := make([]byte, 2)

	// 从连接读数据到 buf 中
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		log.Printf("read connect buf fail , err is %v", err)
		return nil, err
	}

	// validate version
	if buf[0] != Socks5Version {
		return nil, ErrVersionNotSupported
	}

	// read methods
	nmethods := buf[1]
	buf = make([]byte, nmethods)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		log.Printf("read connect nmethods fail , err is %v", err)
		return nil, err
	}

	return &ClientAuthMessage{
		Version:  Socks5Version,
		NMethods: nmethods,
		Methods:  buf,
	}, nil
}

// 认证
func auth(conn net.Conn, conf *Config) error {

	// 读取报文数据
	clientAuthMessage, err := NewClientAuthMessage(conn)
	if err != nil {
		return err
	}
	log.Printf("client auth message is %v", clientAuthMessage)

	// check support method
	var acceptable bool
	for _, method := range clientAuthMessage.Methods {
		if method == conf.AuthMethod {
			acceptable = true
		}
	}

	// 找不到支持的方法
	if !acceptable {
		err = WriteServerAuthMessage(conn, MethodNoAcceptable)
		if err != nil {
			fmt.Printf("server auth message fail,err is %v", err)
			return err
		}
		return ErrMethodNotSupport
	}

	// 找到了可以用的方法
	err = WriteServerAuthMessage(conn, conf.AuthMethod)
	if err != nil {
		return err
	}

	switch conf.AuthMethod {
	case MethodNoAuth:
	case MethodPassWord:
		passwordMessage, pErr := NewClientPasswordMessage(conn)
		if pErr != nil {
			return pErr
		}
		if !conf.PasswordChecker(passwordMessage.Username, passwordMessage.Password) {
			WriteServerPasswordMessage(conn, PasswordAuthFailure)
			return ErrPasswordVerificationFailed
		}
		if pErr = WriteServerPasswordMessage(conn, PasswordAuthSuccess); pErr != nil {
			return pErr
		}
	default:
		return ErrMethodNotSupport
	}
	return nil
}

// WriteServerAuthMessage 发送确认报文
func WriteServerAuthMessage(conn io.Writer, method Method) error {
	buf := []byte{Socks5Version, method}
	_, err := conn.Write(buf)
	if err != nil {
		fmt.Printf("write server auth message fail,err is %v", err)
		return err
	}
	return nil
}

func WriteServerPasswordMessage(conn io.Writer, status byte) error {
	if _, err := conn.Write([]byte{PasswordMethodVersion, status}); err != nil {
		return err
	}
	return nil
}

func NewClientPasswordMessage(conn io.Reader) (*ClientPasswordMessage, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}

	version, usernameLen := buf[0], buf[1]
	if version != PasswordMethodVersion {
		return nil, ErrPasswordVersionNotSupported
	}

	// read username
	buf = make([]byte, usernameLen+1)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}

	username, passwordLen := string(buf[:len(buf)-1]), buf[len(buf)-1]

	// read password
	buf = make([]byte, passwordLen)
	if _, err := io.ReadFull(conn, buf[:passwordLen]); err != nil {
		return nil, err
	}

	return &ClientPasswordMessage{
		Username: username,
		Password: string(buf[:passwordLen]),
	}, nil
}
