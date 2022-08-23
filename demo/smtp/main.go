package main

import (
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/oldbai555/log"
	"gopkg.in/gomail.v2"
	"net/smtp"
)

func main() {
	err := Send4GoMail()
	if err != nil {
		log.Errorf("err:%v", err)
	}
}
func Send4163Email() {
	e := email.NewEmail()
	// 实际发件人
	e.From = "客户 <jiangwenjun@qq.com>"
	// 代发
	e.Sender = "oldbai1005777562@163.com"

	e.To = []string{"wenjunjiang1993@163.com"}

	e.Subject = "学校通知123"

	e.Text = []byte("您好，欢迎入学123")

	err := e.Send("smtp.163.com:25", smtp.PlainAuth("", "oldbai1005777562@163.com", "password", "smtp.163.com"))
	if err != nil {
		log.Errorf("error sending email: %v", err)
	}
}

func Send4GoMail() error {
	// 设置邮箱主体
	mailConn := map[string]string{
		"user": "oldbai1005777562@163.com", // 发送人邮箱（邮箱以自己的为准）
		"pass": "password",                 // 发送人邮箱的密码，现在可能会需要邮箱 开启授权密码后在pass填写授权码
		"host": "smtp.163.com",             // 邮箱服务器（此时用的是qq邮箱）
	}

	m := gomail.NewMessage(
		// 发送文本时设置编码，防止乱码。 如果txt文本设置了之后还是乱码，那可以将原txt文本在保存时
		// 就选择utf-8格式保存
		gomail.SetEncoding(gomail.Base64),
	)

	m.SetHeader("From", "聊天插件 <wenjunjiang1993@163.com>") // 设置发件人
	m.SetHeader("Sender", "oldbai1005777562@163.com")     // 设置实际发件人（用于配制代发）
	m.SetHeader("To", "1005777562@qq.com")                // 发送给用户(可以多个)
	m.SetHeader("Subject", "学校通知12")                      // 设置邮件主题
	m.SetBody("text/html", "您好，欢迎入学1233456")              // 设置邮件正文

	/*
	   创建SMTP客户端，连接到远程的邮件服务器，需要指定服务器地址、端口号、用户名、密码，如果端口号为465的话，
	   自动开启SSL，这个时候需要指定TLSConfig
	*/
	d := gomail.NewDialer(mailConn["host"], 465, mailConn["user"], mailConn["pass"]) // 设置邮件正文
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	dial, err := d.Dial()

	defer func() {
		closeErr := dial.Close()
		if closeErr != nil {
			panic(any(fmt.Sprintf("err:%v", closeErr)))
		}
	}()
	err = gomail.Send(dial, m)
	if err != nil {
		log.Errorf("err:%v", err)
	}
	return err
}
