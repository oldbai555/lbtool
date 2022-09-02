package mail

import (
	"errors"
	"fmt"
)

var connectMap = make(map[uint32]*connect)

const (
	ConnectTypeQQ   uint32 = 1
	ConnectTypeTxEx uint32 = 2
)

func init() {
	connectMap[ConnectTypeQQ] = &connect{SmtpHost: "mailsmtp.qq.com", SmtpPost: 465, ImapHost: "imap.qq.com", ImapPost: 993}
	connectMap[ConnectTypeTxEx] = &connect{SmtpHost: "mailsmtp.exmail.qq.com", SmtpPost: 465, ImapHost: "imap.exmail.qq.com", ImapPost: 993}
}

func GetConnect(connectType uint32) (*connect, error) {
	val, ok := connectMap[connectType]
	if !ok {
		return nil, errors.New(fmt.Sprintf("unknown connect type %d", connectType))
	}
	return val, nil
}

// connect 连接信息
type connect struct {
	SmtpHost string `json:"smtp_host"`
	SmtpPost uint32 `json:"smtp_post"`

	ImapHost string `json:"imap_host"`
	ImapPost uint32 `json:"imap_post"`
}

func NewConnect(smtpHost string, smtpPost uint32, imapHost string, imapPost uint32) *connect {
	return &connect{SmtpHost: smtpHost, SmtpPost: smtpPost, ImapHost: imapHost, ImapPost: imapPost}
}
