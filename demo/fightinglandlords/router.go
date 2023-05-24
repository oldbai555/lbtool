package fightinglandlords

import (
	"encoding/json"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/baix/baix"
	"github.com/oldbai555/lbtool/pkg/baix/iface"
)

const (
	RouterByConnect   uint32 = 1
	RouterByGame      uint32 = 2
	RouterByUnConnect uint32 = 3
)

const (
	ActionJoin      uint32 = 1 // 加入房间
	ActionNewDesk   uint32 = 2 // 创建房间
	ActionGameStart uint32 = 3 // 游戏开始
	ActionPushPoker uint32 = 4 // 出牌
)

var RouterMap = map[uint32]iface.IRouter{
	RouterByConnect:   &ConnectRouter{},
	RouterByGame:      &GameRouter{},
	RouterByUnConnect: &UnConnectRouter{},
}

// ConnectRouter 连接路由
type ConnectRouter struct {
	baix.BaseRouter
}

// GameRouter 游戏路由
type GameRouter struct {
	baix.BaseRouter
}

// UnConnectRouter 断开连接
type UnConnectRouter struct {
	baix.BaseRouter
}

type ConnectRouterReq struct {
	Name string `json:"name"`
}

type ConnectRouterRsp struct {
	Player *Player `json:"player"`
}

type GameRouterReq struct {
	Action uint32 `json:"action"` // 行为

	DeskId string `json:"desk_id"` // 房间ID

	PokerList PokerList `json:"poker_list"`
}

type GameRouterRsp struct{}

type UnConnectRouterReq struct{}

type UnConnectRouterRsp struct{}

func (c *ConnectRouter) Handle(req iface.IRequest) error {
	var r ConnectRouterReq
	err := json.Unmarshal(req.GetData(), &r)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	player := NewPlayer(req.GetConn().GetConnId(), r.Name)

	err = req.Write(&ConnectRouterRsp{
		Player: player,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (c *GameRouter) Handle(req iface.IRequest) error {
	var r GameRouterReq
	err := json.Unmarshal(req.GetData(), &r)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = doAction(req, &r)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = req.Write(&GameRouterRsp{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (c *UnConnectRouter) Handle(req iface.IRequest) error {
	var r UnConnectRouterReq
	err := json.Unmarshal(req.GetData(), &r)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// do_logic

	err = req.Write(&UnConnectRouterRsp{})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

// 执行动作
func doAction(req iface.IRequest, r *GameRouterReq) error {
	player, err := pyMgr.Get(req.GetConn().GetConnId())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	switch r.Action {
	case ActionJoin:
		desk, err := dkMgr.Get(r.DeskId)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = player.joinDesk(desk)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	case ActionNewDesk:
		desk, err := player.NewDesk()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		log.Infof("desk:%v", desk)

		// 需要通知所有人
		header := baix.NewHeader()
		header.SetMsg(req.GetMsg())
		header.SetTraceId(req.GetTraceId())
		payload := baix.NewPayload()
		payload.SetData([]byte(fmt.Sprintf("create a new desk: %v", desk)))
		err = req.GetConn().GetNotify().NotifyBuffAll(baix.NewMessage(header, payload))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	case ActionGameStart:
		desk, err := dkMgr.Get(player.DeskId)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		if !player.IsDeskOwner {
			return fmt.Errorf("not desk owner , player is %v", player)
		}

		// 开始游戏
		err = desk.Start()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// 需要通知所有人这个房间开始游戏了
		header := baix.NewHeader()
		header.SetMsg(req.GetMsg())
		header.SetTraceId(req.GetTraceId())
		payload := baix.NewPayload()
		payload.SetData([]byte(fmt.Sprintf("start game: %v , current player is %v", desk, desk.CurrentPlayer)))
		err = req.GetConn().GetNotify().NotifyBuffAll(baix.NewMessage(header, payload))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

	case ActionPushPoker:
		desk, err := dkMgr.Get(player.DeskId)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		if desk.TurnEnd {
			log.Errorf("current desk %v is end", desk)
			return fmt.Errorf("current desk %v is end", desk)
		}

		if desk.CurrentPlayer.Id != player.Id {
			log.Errorf("current player %s not player %s", desk.CurrentPlayer.Id, player.Id)
			return fmt.Errorf("current player %s not player %s", desk.CurrentPlayer.Id, player.Id)
		}
		defer func() {
			desk.TurnNext() // 轮到下一个人出牌
		}()

		pushPokers := player.PushPokers(r.PokerList)

		header := baix.NewHeader()
		header.SetMsg(req.GetMsg())
		header.SetTraceId(req.GetTraceId())
		payload := baix.NewPayload()
		var data []byte
		if len(pushPokers) > 0 {
			// 通知当前房间的人 出的什么牌
			data = []byte(fmt.Sprintf("current player is %v ,push pokers is %v", desk.CurrentPlayer, pushPokers))
		} else {
			// 通知当前房间的人 要不起
			data = []byte(fmt.Sprintf("current player is %v ,要不起", desk.CurrentPlayer))
		}
		payload.SetData(data)
		err = req.GetConn().GetNotify().NotifyBuffAll(baix.NewMessage(header, payload))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// 出的牌加入牌桌里
		desk.PlayPokers = append(desk.PlayPokers, pushPokers...)

		// 检查手牌出完了没有
		if player.CheckHandPokersEnd() {
			var landlordName = "农民"
			if player.IsLandlord {
				landlordName = "地主"
			}
			// 出完了 游戏结束
			payload.SetData([]byte(fmt.Sprintf("%s【%s】获胜", landlordName, player.Name)))
			err = req.GetConn().GetNotify().NotifyBuffAll(baix.NewMessage(header, payload))
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}

		// 游戏结束
		desk.TurnEnd = true
	}

	return nil
}
