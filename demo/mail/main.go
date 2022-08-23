package main

import (
	"github.com/oldbai555/lb/comm/mail"
	"github.com/oldbai555/lb/log"
)

func main() {
	SendCloudDemo()
}

// SendCloudDemo sendcloud 配制
func SendCloudDemo() {
	err := mail.SendMail(&mail.Sender{
		AuthEmail: "api_user",
		AuthCode:  "api_key",
		SmtpHost:  "smtp.sendcloud.net",
		SmtpPort:  587,
	}, &mail.Details{
		Form:        "username",
		Alias:       "liheng",
		ContentType: mail.DefaultContentType,
		Subject:     "这是一封离别信",
		Body:        []byte("写了许多的消息"),
		ToList:      []string{"username1"},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
