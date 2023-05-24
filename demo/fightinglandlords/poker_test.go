package fightinglandlords

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"testing"
	"time"
)

func TestInitDeskPlayCards(t *testing.T) {
	owner := NewPlayer(utils.GenUUID(), "房主")

	desk, err := owner.NewDesk()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	owner.Ready()

	player := NewPlayer(utils.GenUUID(), "")
	err = player.JoinDesk(desk)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	player.Ready()

	player = NewPlayer(utils.GenUUID(), "")
	err = player.JoinDesk(desk)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	player.Ready()

	for i := 0; i < 16; i++ {
		err = desk.Start()
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}
	}

	time.Sleep(time.Minute)
}
