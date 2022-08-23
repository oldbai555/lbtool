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
	connectMap[ConnectTypeQQ] = &connect{SmtpHost: "smtp.qq.com", SmtpPost: 465, ImapHost: "imap.qq.com", ImapPost: 993}
	connectMap[ConnectTypeTxEx] = &connect{SmtpHost: "smtp.exmail.qq.com", SmtpPost: 465, ImapHost: "imap.exmail.qq.com", ImapPost: 993}
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

// user 收信用户
type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(username string, password string) *user {
	return &user{
		Username: username,
		Password: password,
	}
}

// ReceiveMailMessage 收信内容
type ReceiveMailMessage struct {
	// The message sequence number. It must be greater than or equal to 1.
	SeqNum uint32 `json:"seq_num"`
	// The message envelope.
	Envelope *Envelope `json:"envelope"`
	// The message size.
	Size uint32 `json:"size"`
	// @desc: This is the message's text (can be plain-text or HTML)
	MsgBodyList [][]byte `json:"msg_body_list"`
	// @desc: This is an attachment
	AttachList []*Attach `json:"attach_list"`
	// @desc: The Message-Id header.
	MessageId string `json:"message_id"`
}

// Envelope 信封结构
type Envelope struct {
	// @desc: The message date.
	Date uint32 `json:"date"`
	// @desc: The message subject.
	Subject string `json:"subject"`
	// @desc: The From header addresses.
	FromList []*Address `json:"from_list"`
	// @desc: The message senders.
	SenderList []*Address `json:"sender_list"`
	// @desc: The Reply-To header addresses.
	ReplyToList []*Address `json:"reply_to_list"`
	// @desc: The To header addresses.
	ToList []*Address `json:"to_list"`
	// @desc: The Cc header addresses.
	CcList []*Address `json:"cc_list"`
	// @desc: The Bcc header addresses.
	BccList []*Address `json:"bcc_list"`
	// @desc: The In-Reply-To header. Contains the parent Message-Id.
	InReplyTo string `json:"in_reply_to"`
	// @desc: The Message-Id header.
	MessageId string `json:"message_id"`
}

// Address 信封地址信息
type Address struct {
	// @desc: The personal name.
	PersonalName string `json:"personal_name"`
	// @desc: The SMTP at-domain-list (source route).
	AtDomainList string `json:"at_domain_list"`
	// @desc: The mailbox name.
	MailboxName string `json:"mailbox_name"`
	// @desc: The host name.
	HostName string `json:"host_name"`
}

// Attach 邮箱附件
type Attach struct {
	Buf      []byte `json:"buf"`
	FileName string `json:"file_name"`
}
