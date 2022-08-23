package main

import (
	"github.com/oldbai555/comm/mail"
	"github.com/oldbai555/log"
)

func main() {
	err := mail.SendMail(&mail.Sender{
		AuthEmail: "oldbai1005777562@163.com",
		AuthCode:  "OOHQRFSRSNTJJVOK",
		SmtpHost:  "smtp.163.com",
		SmtpPort:  465,
	}, &mail.Details{
		Form:        "q346407440@gmail.com",
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
