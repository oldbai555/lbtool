package fightinglandlords

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"math/rand"
	"sync"
	"time"
)

const (
	MaxPlayers           = 3 // 最多三人游戏
	MacCallLandlordRound = 4 // 最多叫四轮地主
)

// Desk 牌桌
type Desk struct {
	Id string

	Players    []*Player // 玩家
	Pokers     PokerList // 扑克牌
	PlayPokers PokerList // 出牌列表
	lock       sync.Mutex

	CurrentPlayer *Player // 当前出牌对象
	TurnEnd       bool    // 轮次结束
}

func NewDesk() *Desk {
	d := &Desk{
		Id: utils.GenUUID(),
	}
	dkMgr.Add(d)
	return d
}

// 房间加入玩家
func (d *Desk) joinPlayer(player *Player) error {
	d.lock.Lock()
	defer func() {
		d.lock.Unlock()
		d.initPushSequence()
	}()

	if err := d.CheckAbleToJoinPlayer(); err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	d.Players = append(d.Players, player)
	return nil
}

// 房间离开玩家
func (d *Desk) leavePlayer(player *Player) error {
	d.lock.Lock()
	defer func() {
		d.lock.Unlock()
		d.initPushSequence()
	}()

	var i = -1
	for index, p := range d.Players {
		if p.Id == player.Id {
			i = index
		}
	}

	if i < 0 {
		return PlayerNotFoundLeave
	}

	// 得到剩下的人
	var newPlayers []*Player
	newPlayers = append(newPlayers, d.Players[:i]...)
	if i+1 < len(d.Players) {
		newPlayers = append(newPlayers, d.Players[i+1:]...)
	}
	d.Players = newPlayers

	return nil
}

// 初始化出牌顺序
func (d *Desk) initPushSequence() {
	for i := range d.Players {
		if len(d.Players) == 1 {
			return
		}
		if i+1 < len(d.Players) {
			d.Players[i].SetNext(d.Players[i+1])
		}
		if i+1 == len(d.Players) {
			d.Players[i].SetNext(d.Players[0])
		}
	}
}

// DismissDesk 解散房间
func (d *Desk) DismissDesk() error {
	for _, player := range d.Players {
		err := player.leaveDesk(d)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}
	// 全局删除桌子
	dkMgr.Del(d)
	return nil
}

// Start 开始游戏
func (d *Desk) Start() error {
	if err := d.CheckAbleToStart(); err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 重置对局
	defer d.ResetGame()

	// 初始化一副牌
	d.Pokers = NewPokers()

	// 地主牌
	landlordPokers := d.Pokers[:MaxLandlordPokers]

	// 展示一下地主牌
	log.Infof("地主牌【%s】", landlordPokers.String())

	// 发牌
	d.LicensePoker(d.Pokers[MaxLandlordPokers:])

	// 展示牌
	d.ShowPokers()

	// 取一个随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 开始叫地主的那个人
	firstCallLandlordIndex := r.Intn(MaxPlayers)

	// 轮流叫地主 最多叫四轮
	for i := firstCallLandlordIndex; i < firstCallLandlordIndex+MacCallLandlordRound; i++ {
		d.Players[switchReceivePlayer(i)].CallLandlord()
	}

	// 拿到叫地主的那个人
	landlordPlayerIndex, err := d.GetLandlordRoundPlayers(firstCallLandlordIndex)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("地主【%d】", landlordPlayerIndex)

	// 得到地主的序号,给它加入地主牌
	d.Players[switchReceivePlayer(landlordPlayerIndex)].InsertLandlordPokers(landlordPokers)

	// 当前出牌对象
	d.CurrentPlayer = d.Players[switchReceivePlayer(landlordPlayerIndex)]

	// 自动出牌
	// autoPushPoker(d, d.CurrentPlayer)

	return nil
}

// CheckAbleToStart 检查是否可以开局
func (d *Desk) CheckAbleToStart() (err error) {
	// 人数不满
	if len(d.Players) < MaxPlayers {
		err = DeskPlayerNotUpperLimit
		return
	}

	for _, player := range d.Players {
		if !player.IsReady {
			return DeskPlayerUnReady
		}
	}

	return
}

// CheckAbleToJoinPlayer 检查是否可以加入房间
func (d *Desk) CheckAbleToJoinPlayer() (err error) {
	// 人数达到上限了
	if len(d.Players) >= MaxPlayers {
		err = DeskPlayerUpperLimit
	}
	return
}

// LicensePoker 发牌
func (d *Desk) LicensePoker(list PokerList) {
	for i, card := range list {
		var player *Player
		player = d.Players[switchReceivePlayer(i)]
		player.HandPokerList = append(player.HandPokerList, card)
	}
}

// ShowPokers 展示牌
func (d *Desk) ShowPokers() {
	log.Infof("总牌桌的牌【%s】", d.Pokers.String())

	log.Infof("总出牌记录【%s】", d.PlayPokers.String())

	// 依次展示玩家的牌
	for _, player := range d.Players {
		player.ShowPokers()
	}
}

// ResetGame 重置对局
func (d *Desk) ResetGame() {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 重置玩家手牌信息
	for _, player := range d.Players {
		player.ResetPoker()
	}

	// 重置牌局信息
	d.Pokers = []*Poker{}
	d.PlayPokers = []*Poker{}
	d.TurnEnd = false
}

// GetLandlordRoundPlayers 初始化地主
func (d *Desk) GetLandlordRoundPlayers(firstPlayerIndex int) (int, error) {
	// 最开始叫的那个人, 有没有在最后叫地主
	player := d.Players[switchReceivePlayer(firstPlayerIndex)]
	if player.AskLandlord {
		// 初始化地主
		player.IsLandlord = true
		d.Players[switchReceivePlayer(firstPlayerIndex)] = player
		return firstPlayerIndex, nil
	}

	// 如果没有,那就递归找它的上一个
	for i := 2; i > 0; i-- {
		player := d.Players[switchReceivePlayer(firstPlayerIndex+i)]
		if player.AskLandlord {
			// 初始化地主
			player.IsLandlord = true
			d.Players[switchReceivePlayer(firstPlayerIndex+i)] = player
			return firstPlayerIndex + i, nil
		}
	}

	// 可以做都没人叫 把第一个人当地主
	log.Errorf("err is %v,选第一个叫地主的人当地主", DeskNotLandlordRound)
	return firstPlayerIndex, nil
}

// TurnNext 轮转到下一个人,并且返回下一个人是谁
func (d *Desk) TurnNext() (next *Player) {
	d.CurrentPlayer = d.CurrentPlayer.Next
	return d.CurrentPlayer
}

func autoPushPoker(d *Desk, player *Player) {
	// 初始第一个出牌的人是地主
	// 开始出牌
	var lashPlayerPushPokers []*Poker
	var lastPushPlayer = player
	for {
		// 判断是否是本人
		if player.Id != d.CurrentPlayer.Id {
			log.Errorf("current player %s not player %s", d.CurrentPlayer.Id, player.Id)
			continue
		}

		log.Infof("轮到【%s】出牌", player.Name)

		// 上一次出牌的用户是自己,表示一圈都没有人要得起牌,需要重置一下上一次出牌
		if lastPushPlayer.Id == player.Id {
			lashPlayerPushPokers = []*Poker{}
		}

		// 出牌
		pokerList := player.AutoPushPokers(lashPlayerPushPokers)
		if len(pokerList) == 0 {
			log.Infof("【%s】,要不起", player.Name)
			// 轮到下一个人
			player = d.TurnNext()
			continue
		}

		log.Infof("出牌为【%s】", pokerList.String())

		// 出的牌加入牌桌里
		d.PlayPokers = append(d.PlayPokers, pokerList...)

		// 记录到上一次出的牌
		lashPlayerPushPokers = pokerList

		// 记录上一次出牌的用户
		lastPushPlayer = player

		// 检查手牌出完了没有
		if player.CheckHandPokersEnd() {
			var landlordName = "农民"
			if player.IsLandlord {
				landlordName = "地主"
			}
			// 出完了 游戏结束
			log.Infof("%s【%s】获胜", landlordName, player.Name)
			break
		}

		// 轮到下一个玩家出牌
		player = d.TurnNext()

		time.Sleep(time.Second)
	}

	log.Infof("总出牌记录【%s】", d.PlayPokers.String())

	// 展示牌
	d.ShowPokers()

	// 重置下一局游戏
	d.ResetGame()
}
