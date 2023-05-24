package proto_emicklei

import (
	"github.com/oldbai555/lbtool/log"
	"os"
)

func Exit(err error) {
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
	}
}

func Must(err error) {
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}
}

func LogIf(b bool, err error) {
	if err != nil && b {
		log.Errorf("err:%v", err)
	}
}
