package sixhecai

import (
	"github.com/oldbai555/lbtool/log"
	"testing"
)

func TestParseNumber(t *testing.T) {
	number, err := ParseNumber("0102030405101.100")
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	log.Infof("number is %s", number)
}

func TestParseText(t *testing.T) {
	text, err := ParseText("猪,狗,牛各50\n" +
		"兔龙各100\n" +
		"010203各150\n" +
		"\n" +
		"\n" +
		"\n" +
		"\n" +
		"11.12.13只各粒200\n" +
		"04-05-06只各粒250\n" +
		"41.42.43只各粒300块\n")
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	user := NewUser("张三")
	err = user.SaveXz(text)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}

	number, err := ParseNumber("0102030405101.100")
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	err = user.SaveXz(number)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	user.Show()

	//user.ShowRecord()
}
