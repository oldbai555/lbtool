package mail

import (
	"encoding/json"
	"lb/log"
	"sync"
	"testing"
	"time"
)

// 腾讯企业邮箱 限频每秒约20
// 个邮 授权码被风控
func TestReadMail(t *testing.T) {
	var sync1 sync.WaitGroup
	// var sync2 sync.WaitGroup
	for i := 0; i < 1; i++ {
		go dologic(&sync1, "zhangjianjun@aquanliang.com", "password")
		// go dologic(&sync2,"1005777562@qq.com","password")
	}
	sync1.Wait()
	// sync2.Wait()
	time.Sleep(5 * time.Second)
	return
}

func dologic(s *sync.WaitGroup, username, password string) {
	s.Add(1)
	defer s.Done()
	c, err := NewImapClient(ConnectTypeQQ)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}

	err = ReadMail(&ReadMailReq{
		C:        c,
		U:        NewUser(username, password),
		Flags:    []string{FlagsAll},
		ReadSize: 1,
	},
		func(msgList []*ReceiveMailMessage) error {
			for _, msg := range msgList {
				marshal, _ := json.Marshal(msg)
				log.Infof("msg = %v ", string(marshal))
			}
			log.Infof("len(msgList) = %d ", len(msgList))
			return nil
		})
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}

}

func TestClient(t *testing.T) {
	// 连接邮件服务器
	c, err := NewImapClient(ConnectTypeQQ)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}

	err = Login(NewUser("1005777562@qq.com", "yjxawxmhqwozbcbg"), c)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}

	err = c.Logout()
	if err != nil {
		log.Errorf("err is %v", err)
	}

	err = Login(NewUser("1005777562@qq.com", "yjxawxmhqwozbcbg"), c)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}

	err = c.Logout()
	if err != nil {
		log.Errorf("err is %v", err)
	}
}
