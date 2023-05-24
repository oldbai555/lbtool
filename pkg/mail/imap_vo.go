package mail

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
	// @desc: The SMTP at-bconf-list (source route).
	AtDomainList string `json:"at_domain_list"`
	// @desc: The mailbox name.
	MailboxName string `json:"mailbox_name"`
	// @desc: The host name.
	HostName string `json:"host_name"`
}
