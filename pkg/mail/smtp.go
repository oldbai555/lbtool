package mail

import (
	"crypto/tls"
	"fmt"
	"github.com/oldbai555/lbtool/extpkg/gomail"
	"github.com/oldbai555/lbtool/log"
	"mime"
	"strings"
)

const (
	DefaultContentType = "text/html"

	HeaderFrom    = "From"     // 发件人
	HeaderSender  = "Sender"   // 发件人
	HeaderTo      = "To"       // 收件人
	HeaderReplyTo = "Reply-to" // 答复
	HeaderSubject = "Subject"  // 主题
	HeaderCc      = "Cc"       // 抄送
	HeaderBcc     = "Bcc"      // 密件抄送
)

// SendMail 发送邮件
// sender 和 from 不同时，发件人为 from , 代发人为 sender ，sender 代理 from 发送邮件
// sender 和 from 相同时，发件人和代发人都相同，相当于自己发送
func SendMail(sender *Sender, detail *Details) error {

	s, err := NewSendClient(sender)
	defer func() {
		closeErr := s.Close()
		if closeErr != nil {
			panic(any(fmt.Sprintf("err:%v", closeErr)))
		}
	}()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	m := gomail.NewMessage(
		gomail.SetEncoding(gomail.Base64),
	)

	// 设置发件人信息
	m.SetHeader(HeaderFrom, m.FormatAddress(detail.Form, detail.Alias))

	// 设置实际的发件人,部分邮件服务商允许 From 和 Sender 不一致,实际使用时使用 From 即可;
	// m.SetHeader(HeaderSender, sender.AuthEmail)

	// 设置收件人信息
	m.SetHeader(HeaderTo, detail.ToList...)

	// 指定回复人
	// m.SetHeader(HeaderReplyTo, "934105499@qq.com")

	// 设置邮箱主题
	m.SetHeader(HeaderSubject, detail.Subject)

	// 抄送对象
	if len(detail.CarbonCopyList) > 0 {
		m.SetHeader(HeaderCc, detail.CarbonCopyList...)
	}

	// 加密抄送对象
	if len(detail.BlindCarbonCopyList) > 0 {
		m.SetHeader(HeaderBcc, detail.BlindCarbonCopyList...)
	}

	// 邮件内容
	if len(detail.Body) > 0 {
		ct := detail.ContentType
		if ct == "" {
			ct = DefaultContentType
		}
		m.SetBody(ct, string(detail.Body))
	} else {
		m.SetBody(DefaultContentType, "")
	}

	// 添加附件
	for _, v := range detail.Attach {
		m.AttachBuffer(v.FileName, v.Buf, gomail.SetHeader(map[string][]string{
			"Content-Disposition": {
				fmt.Sprintf(`attachment; filename="%s"`, mime.BEncoding.Encode("UTF-8", v.FileName)),
			},
		}))
	}

	// 开始发送
	if err = gomail.Send(s, m); err != nil {
		log.Errorf("send mail err: %v", err)
		if strings.Contains(err.Error(), "invalid address") {
			return fmt.Errorf("invalid address")
		}
		return err
	}
	return nil
}

// NewSendClient 声明邮件发送Client
func NewSendClient(sender *Sender) (gomail.SendCloser, error) {
	// 声明连接邮箱服务器对象
	d := gomail.NewDialer(sender.SmtpHost, int(sender.SmtpPort), sender.AuthEmail, sender.AuthCode)

	// 关闭TLS认证设置,为true时，关闭TLS认证，否则默认开启，需要配置证书认证
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	// 配制 SSL
	// d.SSL = true

	// 开始建立客户端
	s, err := d.Dial()
	if err != nil {
		return nil, err
	}

	return s, nil
}
