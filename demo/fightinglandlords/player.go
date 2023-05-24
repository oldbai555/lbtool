package fightinglandlords

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Player struct {
	Id     string
	Name   string
	DeskId string

	IsDeskOwner bool // 是否是房主
	IsReady     bool // 是否准备

	HandPokerList PokerList // 手上的牌
	PlayPokerList PokerList // 打出的牌

	IsLandlord  bool // 是否是地主
	AskLandlord bool // 叫地主
	lock        sync.Mutex

	Next *Player
}

func NewPlayer(id, name string) *Player {
	if name == "" {
		name = fmt.Sprintf("玩家%s", id)
	}
	p := &Player{
		Id:          id,
		Name:        name,
		AskLandlord: true, // 默认都想叫地主

		IsReady:     false,
		IsDeskOwner: false,
		IsLandlord:  false,
	}
	pyMgr.Add(p)
	return p
}

// ResetPoker 重置牌信息
func (p *Player) ResetPoker() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.HandPokerList = PokerList{}
	p.PlayPokerList = PokerList{}

	p.IsLandlord = false // 是否是地主
	p.AskLandlord = true // 默认都想叫地主
	p.IsReady = true     // 默认重置牌局时是准备的
}

// ShowPokers 显示手牌和出的牌
func (p *Player) ShowPokers() {
	p.SortHandPokers()
	log.Infof("【%s】手牌【%s】", p.Name, p.HandPokerList.String())
	log.Infof("【%s】出牌记录【%s】", p.Name, p.PlayPokerList.String())
}

// SortHandPokers 排序手牌
func (p *Player) SortHandPokers() {
	sort.Sort(p.HandPokerList)
}

// CallLandlord 叫地主
func (p *Player) CallLandlord() {
	// 第一次进来是可以选择的
	// 如果不叫地主,那么就跳过
	if !p.AskLandlord {
		return
	}

	// 交给时间来判断是否叫地主
	p.AskLandlord = time.Now().Unix()%2 == 0
}

// NewDesk 创一个新的房间
func (p *Player) NewDesk() (*Desk, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if err := p.CheckAbleToJoinDesk(); err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	d := NewDesk()

	// 房主加入房间
	err := p.joinDesk(d)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	p.IsDeskOwner = true

	return d, nil
}

// JoinDesk 加入房间
func (p *Player) JoinDesk(d *Desk) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.joinDesk(d)
}

// 加入房间逻辑
func (p *Player) joinDesk(d *Desk) error {
	if err := p.CheckAbleToJoinDesk(); err != nil {
		log.Errorf("err is %v", err)
		return err
	}

	// 房间加入玩家
	err := d.joinPlayer(p)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 设置房间ID
	p.DeskId = d.Id
	return nil
}

// CheckAbleToJoinDesk 检查是否可以加入房间
func (p *Player) CheckAbleToJoinDesk() (err error) {
	// 房主 或 已经有房间了, 不支持直接跳过去
	if p.IsDeskOwner || p.DeskId != "" {
		err = NotCanJoinDesk
	}
	return
}

// LeaveDesk 退房
func (p *Player) LeaveDesk() error {
	desk, err := dkMgr.Get(p.DeskId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 房主就直接解散房间
	if p.IsDeskOwner {
		err := desk.DismissDesk()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	}

	// 否则就正常退房
	err = p.leaveDesk(desk)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

// 退房逻辑
func (p *Player) leaveDesk(desk *Desk) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	err := desk.leavePlayer(p)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	p.ResetDesk()
	return nil
}

// ResetDesk 重置牌桌信息
func (p *Player) ResetDesk() {
	p.DeskId = ""
	p.IsDeskOwner = false
	p.IsReady = false
	p.SetNext(nil)
}

// ResetUnReady 重置为未准备
func (p *Player) ResetUnReady() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.IsReady = false
}

// Ready 准备
func (p *Player) Ready() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.IsReady = true
}

// InsertLandlordPokers 加入地主牌
func (p *Player) InsertLandlordPokers(pokers PokerList) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.HandPokerList = append(p.HandPokerList, pokers...)
	// 重新排下序
	p.SortHandPokers()
}

// PushPokers 出牌
func (p *Player) PushPokers(pokerList PokerList) (list PokerList) {

	lastHandsPokerList, abovePokerList := diffPushPokerList(p.HandPokerList, pokerList)

	// 放入剩下的手牌
	p.HandPokerList = lastHandsPokerList

	// 排序一下
	p.SortHandPokers()

	// 放入出牌记录
	p.PlayPokerList = append(p.PlayPokerList, abovePokerList...)

	// 放入结果
	list = abovePokerList

	return
}

func (p *Player) AutoPushPokers(lashPlayerPushPokers PokerList) (list PokerList) {
	// 不是第一次出牌,开始找哪几张可以出
	lastHandsPokerList, abovePokerList := pokersAbove(p.HandPokerList, lashPlayerPushPokers)

	// 上一次没人出牌 所以找不到可以出的牌,随机出一张
	if len(abovePokerList) == 0 && len(lashPlayerPushPokers) == 0 {
		// 取一个随机数
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		// 随便出一张
		var ablePushIndex int
		if len(p.HandPokerList)-1 != 0 {
			ablePushIndex = r.Intn(len(p.HandPokerList))
		}

		// 出牌
		poker := p.HandPokerList[ablePushIndex]

		// 记录出牌
		p.PlayPokerList = append(p.PlayPokerList, poker)

		// 重新放入手牌
		var newPokerList = p.HandPokerList[:ablePushIndex]

		// 判断下是否超出了界限
		if len(p.HandPokerList) > ablePushIndex+1 {
			newPokerList = append(newPokerList, p.HandPokerList[ablePushIndex+1:]...)
		}

		// 重新放入
		p.HandPokerList = newPokerList

		// 排序一下
		p.SortHandPokers()

		// 放入结果
		list = append(list, poker)
		return
	}

	// 要得起
	if len(abovePokerList) != 0 {
		// 放入剩下的手牌
		p.HandPokerList = lastHandsPokerList

		// 排序一下
		p.SortHandPokers()

		// 放入出牌记录
		p.PlayPokerList = append(p.PlayPokerList, abovePokerList...)

		// 放入结果
		list = abovePokerList
		return
	}
	return
}

// CheckHandPokersEnd 检查是否出牌完了
func (p *Player) CheckHandPokersEnd() bool {
	return len(p.HandPokerList) == 0
}

// SetNext 设置下家
func (p *Player) SetNext(player *Player) {
	p.Next = player
}
