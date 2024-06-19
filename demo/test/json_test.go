package test

//
//import (
//	"bytes"
//	"google.golang.org/protobuf/jsonpb"
//	jsoniter "github.com/json-iterator/go"
//	"github.com/oldbai555/lbtool/log"
//	"testing"
//)
//
//func TestJsonIter(t *testing.T) {
//	str := "{\"type\":1,\"match_player\":{\"player_id\":\"3\"}}"
//	var event lbddz.Event
//	err := jsoniter.Unmarshal([]byte(str), &event)
//	if err != nil {
//		log.Errorf("err:%v", err)
//	}
//	var event2 = lbddz.Event{
//		Type: 1,
//		MatchPlayer: &lbddz.MatchPlayer{
//			PlayerId: uint64(9223372036854775807),
//		},
//	}
//	marshal, err := jsoniter.Marshal(event2)
//	if err != nil {
//		log.Errorf("err:%v", err)
//	}
//	log.Infof("event:%v", event)
//	log.Infof("event2:%v", string(marshal))
//}
//
//func TestJsonPb(t *testing.T) {
//	str := "{\"type\":1,\"match_player\":{\"player_id\":\"3\"}}"
//	var event lbddz.Event
//	unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: true}
//	err := unmarshaler.Unmarshal(bytes.NewReader([]byte(str)), &event)
//	if err != nil {
//		log.Warnf("invalid json  %v", err)
//		return
//	}
//	log.Infof("event is %v", event)
//
//	var m = jsonpb.Marshaler{
//		EmitDefaults: true,
//		OrigName:     true,
//	}
//	var buf []byte
//	var temp bytes.Buffer
//	err = m.Marshal(&temp, &lbddz.Event{
//		Type: 1,
//		MatchPlayer: &lbddz.MatchPlayer{
//			PlayerId: uint64(9223372036854775807),
//		},
//	})
//	if err != nil {
//		log.Errorf("proto MarshalToString err:", err)
//		return
//	}
//	buf = temp.Bytes()
//	log.Infof("buf is %s", buf)
//}
