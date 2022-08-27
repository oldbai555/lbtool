package mail

import (
	"github.com/oldbai555/lb/log"
	"testing"
)

func TestSendMail(t *testing.T) {
	err := SendMail(&Sender{
		AuthEmail: "api_user",
		AuthCode:  "api_key",
		SmtpHost:  "smtp.sendcloud.net",
		SmtpPort:  587,
	}, &Details{
		Form:        "username",
		Alias:       "liheng",
		ContentType: DefaultContentType,
		Subject:     "这是一封离别信",
		Body:        []byte("写了许多的消息"),
		ToList:      []string{"username1"},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
