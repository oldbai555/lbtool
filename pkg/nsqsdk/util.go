package nsqsdk

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"time"
)

const (
	TIMEOUT = 5 * time.Second
)

type Msg struct {
	ReqId   string
	CorpId  uint32
	MsgInfo []byte
}

func EncodeMsg(reqId string, corpId uint32, msg []byte) ([]byte, error) {
	info := &Msg{
		ReqId:   reqId,
		CorpId:  corpId,
		MsgInfo: msg,
	}
	b, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func DecodeMsg(msg []byte) (*Msg, error) {
	m := new(Msg)
	err := json.Unmarshal(msg, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Process(msg *nsq.Message, doLogic func(data interface{}) error) error {
	info, err := DecodeMsg(msg.Body)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return err
	}

	if msg.Attempts > 3 {
		fmt.Printf("Exceeding maximum limit %v", string(info.MsgInfo))
		msg.Finish()
		return nil
	}

	var data interface{}
	err = json.Unmarshal(info.MsgInfo, &data)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return err
	}

	err = doLogic(data)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return err
	}
	return nil
}
