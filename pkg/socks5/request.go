package socks5

import (
	"io"
	"net"
)

const (
	IPv4Length = 4
	Ipv6Length = 16
	PortLength = 2
)

type ClientRequestMessage struct {
	Version  byte        `json:"version"`
	Cmd      Command     `json:"cmd"`
	Reserved byte        `json:"reserved"`
	Atyp     AddressType `json:"atyp"`
	Address  string      `json:"address"`
	Port     uint16      `json:"port"`
}

type Command = byte

const (
	CmdConnect Command = 0x01
	CmdBind    Command = 0x02
	CmdUDP     Command = 0x03
)

type AddressType = byte

const (
	AddressTypeIPv4   AddressType = 0x01
	AddressTypeDomain AddressType = 0x03
	AddressTypeIPv6   AddressType = 0x04
)

type ReplyType = byte

const (
	ReplyTypeSuccess ReplyType = iota
	ReplyTypeServerFailure
	ReplyTypeConnectionNotAllowed
	ReplyTypeNewWorkUnreachable
	ReplyTypeHostUnreachable
	ReplyTypeConnectionRefused
	ReplyTypeTTLExpired
	ReplyTypeCommandNotSupported
	ReplyTypeAddressTypeNotSupported
)

func NewClientRequestMessage(conn io.Reader) (*ClientRequestMessage, error) {
	var crm = ClientRequestMessage{}
	// read version, command, reserved, address type
	var buf = make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}

	crm.Version, crm.Cmd, crm.Reserved, crm.Atyp = buf[0], buf[1], buf[2], buf[3]

	// validate version
	if crm.Version != Socks5Version {
		return nil, ErrVersionNotSupported
	}

	switch crm.Cmd {
	case CmdConnect, CmdBind, CmdUDP:
	default:
		return nil, ErrCommandNotSupported
	}

	if crm.Reserved != ReservedField {
		return nil, ErrReservedFieldNotSupported
	}

	switch crm.Atyp {
	case AddressTypeIPv6:
		buf = make([]byte, Ipv6Length)
		fallthrough
	case AddressTypeIPv4:
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
		ip := net.IP(buf)
		crm.Address = ip.String()
		break
	case AddressTypeDomain:
		if _, err := io.ReadFull(conn, buf[:1]); err != nil {
			return nil, err
		}
		domainLength := buf[0]
		if domainLength > IPv4Length {
			buf = make([]byte, domainLength)
		}
		if _, err := io.ReadFull(conn, buf[:domainLength]); err != nil {
			return nil, err
		}
		crm.Address = string(buf[:domainLength])
	default:
		return nil, ErrAddressTypeNotSupported
	}

	if _, err := io.ReadFull(conn, buf[:PortLength]); err != nil {
		return nil, err
	}
	crm.Port = uint16(buf[0]<<8) + uint16(buf[1])

	return &crm, nil
}

func WriteRequestSuccessMessage(conn io.Writer, ip net.IP, port uint16) error {
	var addressType = AddressTypeIPv4
	if len(ip) == Ipv6Length {
		addressType = AddressTypeIPv6
	}

	// write version ,reply success,reserved , address type
	if _, err := conn.Write([]byte{
		Socks5Version,
		ReplyTypeSuccess,
		ReservedField,
		addressType,
	}); err != nil {
		return err
	}

	// write bind Ip (ipv4/ipv6)
	if _, err := conn.Write(ip); err != nil {
		return err
	}

	// write bind port
	buf := make([]byte, 2)
	buf[0] = byte(port >> 8)
	buf[1] = byte(port - uint16(buf[0])<<8)
	if _, err := conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func WriteRequestFailureMessage(conn io.Writer, replyType ReplyType) error {
	if _, err := conn.Write([]byte{
		Socks5Version,
		replyType,
		ReservedField,
		0, 0, 0, 0,
		0, 0,
	}); err != nil {
		return err
	}
	return nil
}

// TcpForward 请求转发
func TcpForward(conn io.ReadWriter, targetConn io.ReadWriteCloser) error {
	defer targetConn.Close()
	go io.Copy(targetConn, conn)
	_, err := io.Copy(conn, targetConn)
	return err
}
