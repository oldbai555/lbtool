package main

import (
	"github.com/oldbai555/lb/comm/mail"
	"github.com/oldbai555/lb/log"
)

func main() {
	SendCloudDemo()
}

// SendCloudDemo sendcloud 配制
// sc_rtxkw0_test_o0szgV api_user
// b4d22b6261626eccc359254ca8530a53 api_code
// 9g6OlrWtRaKKt8vqL8FTWugxKsoK9Rel.sendcloud.org send_host
func SendCloudDemo() {
	err := mail.SendMail(&mail.Sender{
		AuthEmail: "sc_rtxkw0_test_o0szgV",
		AuthCode:  "b4d22b6261626eccc359254ca8530a53",
		SmtpHost:  "smtp.sendcloud.net",
		SmtpPort:  587,
	}, &mail.Details{
		Form:        "wenjunjiang1993@163.com",
		Alias:       "liheng",
		ContentType: mail.DefaultContentType,
		Subject:     "这是一封离别信",
		Body:        []byte("写了许多的消息"),
		ToList:      []string{"1005777562@qq.com"},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
