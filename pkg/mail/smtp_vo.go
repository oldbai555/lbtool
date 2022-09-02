package mail

// Sender 发送者信息
type Sender struct {
	// AuthEmail 发件人邮件
	AuthEmail string `json:"send_mail"`

	// AuthCode 授权码
	AuthCode string `json:"password"`

	// SmtpHost smtp域名
	SmtpHost string `json:"smtp_host"`

	// SmtpPort smtp端口
	SmtpPort uint32 `json:"smtp_port"`
}

// Details 邮件内容
type Details struct {
	// Form 发件人
	Form string `json:"form_mail"`

	// Alias 发件人别名
	Alias string `json:"alias"`

	// ContentType default html
	ContentType string `json:"content_type"`

	// Subject 主题
	Subject string `json:"subject"`

	// Body 发送内容
	Body []byte `json:"body"`

	// Attach 邮件附件
	Attach []*Attach `json:"attach"`

	// CarbonCopyList 抄送列表
	CarbonCopyList []string `json:"carbon_copy_list"`

	// BlindCarbonCopyList 密送列表
	BlindCarbonCopyList []string `json:"blind_carbon_copy_list"`

	// ToList 收件人列表
	ToList []string `json:"to_mail_list"`
}

// Attach 邮件附件
type Attach struct {
	Buf      []byte `json:"buf"`
	FileName string `json:"file_name"`
}
