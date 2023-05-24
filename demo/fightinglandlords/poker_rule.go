package fightinglandlords

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"sort"
	"strings"
)

var (
	// Pokers 牌型 -> 组合
	Pokers = make(map[string]*Combination, 16384)
	// TypeToPokers 牌型 -> 组合列表
	TypeToPokers = make(map[string][]*Combination, 38)
)

// Combination 扑克组合规则
type Combination struct {
	Type   string `json:"type"`   // 牌型
	Score  int    `json:"score"`  // 牌分 - 也就是牌的大小
	Pokers string `json:"pokers"` // 牌的字符串
}

var pokersTypeMap = map[string]string{
	SeqSingle7:     "七张顺子",
	SeqPair6:       "六张顺子",
	BombSingle:     "四带两单牌",
	SeqTrioPair2:   "两飞机带两对子",
	SeqTrioPair5:   "五飞机带五对子",
	SeqPair7:       "七连对",
	TrioSingle:     "三带一",
	SeqSingle11:    "11张顺子",
	SeqTrio6:       "六飞机不带",
	SeqPair8:       "八连对",
	SeqPair3:       "三连对",
	SeqTrioPair4:   "四飞机带四对",
	SeqSingle10:    "十张顺子",
	SeqSingle9:     "九张顺子",
	SeqSingle5:     "五张顺子",
	SeqTrio2:       "两飞机不带",
	Pair:           "一对子",
	Trio:           "三连张不带",
	SeqPair5:       "五连对",
	Bomb:           "炸弹",
	SeqTrio5:       "五飞机不带",
	SeqSingle12:    "十二张顺子",
	SeqTrioPair3:   "三飞机带三对",
	BombPair:       "四带两对",
	SeqTrioSingle4: "四飞机带四张单牌",
	SeqTrioSingle5: "五飞机带五张单牌",
	SeqPair10:      "十连对",
	SeqSingle8:     "八张顺子",
	SeqTrioSingle3: "三飞机带三张单牌",
	SeqTrio4:       "四飞机不带",
	SeqPair4:       "四连对",
	Rocket:         "王炸",
	SeqPair9:       "九连对",
	SeqSingle6:     "六张顺子",
	SeqTrioSingle2: "飞机带两张单牌",
	SeqTrio3:       "三飞机不带",
	TrioPair:       "三带一对子",
	Single:         "单牌",
}

const (
	SeqSingle7     = "seq_single7"      // 七张顺子
	SeqPair6       = "seq_pair6"        // 六张顺子
	BombSingle     = "bomb_single"      // 四带两单牌
	SeqTrioPair2   = "seq_trio_pair2"   // 两飞机带两对子
	SeqTrioPair5   = "seq_trio_pair5"   // 五飞机带五对子
	SeqPair7       = "seq_pair7"        // 七连对
	TrioSingle     = "trio_single"      // 三带一
	SeqSingle11    = "seq_single11"     // 11张顺子
	SeqTrio6       = "seq_trio6"        // 六飞机不带
	SeqPair8       = "seq_pair8"        // 八连对
	SeqPair3       = "seq_pair3"        // 三连对
	SeqTrioPair4   = "seq_trio_pair4"   // 四飞机带四对
	SeqSingle10    = "seq_single10"     // 十张顺子
	SeqSingle9     = "seq_single9"      // 九张顺子
	SeqSingle5     = "seq_single5"      // 五张顺子
	SeqTrio2       = "seq_trio2"        // 两飞机不带
	Pair           = "pair"             // 一对子
	Trio           = "trio"             // 三连张不带
	SeqPair5       = "seq_pair5"        // 五连对
	Bomb           = "bomb"             // 炸弹
	SeqTrio5       = "seq_trio5"        // 五飞机不带
	SeqSingle12    = "seq_single12"     // 十二张顺子
	SeqTrioPair3   = "seq_trio_pair3"   // 三飞机带三对
	BombPair       = "bomb_pair"        // 四带两对
	SeqTrioSingle4 = "seq_trio_single4" // 四飞机带四张单牌
	SeqTrioSingle5 = "seq_trio_single5" // 五飞机带五张单牌
	SeqPair10      = "seq_pair10"       // 十连对
	SeqSingle8     = "seq_single8"      // 八张顺子
	SeqTrioSingle3 = "seq_trio_single3" // 三飞机带三张单牌
	SeqTrio4       = "seq_trio4"        // 四飞机不带
	SeqPair4       = "seq_pair4"        // 四连对
	Rocket         = "rocket"           // 王炸
	SeqPair9       = "seq_pair9"        // 九连对
	SeqSingle6     = "seq_single6"      // 六张顺子
	SeqTrioSingle2 = "seq_trio_single2" // 飞机带两张单牌
	SeqTrio3       = "seq_trio3"        // 三飞机不带
	TrioPair       = "trio_pair"        // 三带一对子
	Single         = "single"           // 单牌
)

//go:embed rule.json
var jsonStrByte string

func init() {
	var rule = make(map[string][]string)
	err := json.Unmarshal([]byte(jsonStrByte), &rule)
	if err != nil {
		panic(fmt.Sprintf("json unmarsha err:%v", err))
		return
	}

	// pokerType 牌型 , pokers 牌的组合列表, 已经按从小到大排序
	for pokerType, pokers := range rule {

		// score 用下标来表示牌的大小 poker 具体的牌组合
		for score, poker := range pokers {
			pokers := SortStr(poker)
			p := &Combination{
				Type:   pokerType,
				Score:  score,
				Pokers: pokers,
			}
			// 指定这个牌型是这个组合
			Pokers[pokers] = p

			// 牌型组合列表
			TypeToPokers[pokerType] = append(TypeToPokers[pokerType], p)
		}
	}
}

// 手牌是否存在指定牌型
func pokersInContains(parent, child string) bool {
	return strings.Contains(parent, child)
}

// 获得牌型和大小
func pokersValue(pokers string) (cardType string, score int) {
	if combination, ok := Pokers[SortStr(pokers)]; ok {
		cardType = combination.Type
		score = combination.Score
	}
	return
}

// 比较牌大小,并返回是否翻倍
// lastShotPoker 上一次打的牌
// comparedNum 这一次打的牌
// -1 上一次打的牌 小于 这一次打的牌
// 0 上一次打的牌 相等 这一次打的牌
// 1 上一次打的牌 大于 这一次打的牌
func comparePokers(lastShotPokers, comparedNum PokerList) (compareRes int, isMulti bool) {
	log.Infof("comparedNum %v  %v", lastShotPokers, comparedNum)

	// 如果没有牌了 或者比较的牌没有了
	if len(lastShotPokers) == 0 || len(comparedNum) == 0 {
		// 进来了说明肯定有一个不出
		if len(lastShotPokers) == 0 && len(comparedNum) == 0 {
			// 都没有出牌
			return
		} else {
			// 上一次打的牌不为空
			if len(lastShotPokers) != 0 {
				// 上一次打的牌不为空,那这次打的牌肯定为空
				compareRes = -1
				return
			} else {
				// 上一次打的牌为空,那这次打的牌肯定比上一次大
				compareRes = 1

				// 获取比较的牌型
				comparedType, _ := pokersValue(comparedNum.String())

				// 炸弹翻倍
				if comparedType == Rocket || comparedType == Bomb {
					isMulti = true
					return
				}
				return
			}
		}
	}

	// 得到牌型和大小
	lastShotPokerType, lastShotPokerScore := pokersValue(lastShotPokers.String())

	// 得到牌型和大小
	comparedType, comparedScore := pokersValue(comparedNum.String())

	log.Infof("compare poker %v, %v, %v, %v", lastShotPokerType, lastShotPokerScore, comparedType, comparedScore)

	// 牌型相同,直接比大小
	if lastShotPokerType == comparedType {
		compareRes = comparedScore - lastShotPokerScore
		return
	}

	// 这次打的是王炸
	if comparedType == Rocket {
		isMulti = true
		compareRes = 1
		return
	}

	// 上一次打的是王炸
	if lastShotPokerType == Rocket {
		compareRes = -1
		return
	}

	// 这次打的是普通的炸弹
	if comparedType == Bomb {
		isMulti = true
		compareRes = 1
		return
	}

	return
}

// 查找手牌中是否有比被比较牌型大的牌
// lastHandsNum 剩下的手牌
// aboveNum 可以打的牌
func pokersAbove(handsPokerList, lastPushPokerList PokerList) (lastHandsPokerList, abovePokerList PokerList) {
	// 手牌
	handCards := handsPokerList.String()

	// 别人出的牌
	turnCards := lastPushPokerList.String()

	// 拿到别人的牌型和牌型大小
	cardType, cardScore := pokersValue(turnCards)

	log.Infof("比较牌 ==> 手牌【%s】,牌型【%s】,上一次出的牌【%s】", handCards, pokersTypeMap[cardType], turnCards)
	// 找不到牌型
	if len(cardType) == 0 {
		lastHandsPokerList = handsPokerList
		return
	}

	// 去找这个牌型中比它大的牌
	for _, combination := range TypeToPokers[cardType] {
		// 比它大 且 存在
		if combination.Score > cardScore && pokersInContains(handCards, combination.Pokers) {
			// 转换为可以出的牌
			lastHandsPokerList, abovePokerList = convertPushPokerList(handsPokerList, combination.Pokers)
			return
		}
	}

	// 找不到就去找炸弹
	if cardType != Bomb && cardType != Rocket {
		// 开始找炸弹
		for _, combination := range TypeToPokers[Bomb] {
			// 手牌里有符合的炸弹
			if pokersInContains(handCards, combination.Pokers) {
				// 转换为可以出的牌
				lastHandsPokerList, abovePokerList = convertPushPokerList(handsPokerList, combination.Pokers)
				return
			}
		}
	} else if pokersInContains(handCards, KingRocket) {
		// 转换为可以出的牌
		lastHandsPokerList, abovePokerList = convertPushPokerList(handsPokerList, KingRocket)
		return
	}

	// 走到这里表示没有可以出的牌
	lastHandsPokerList = handsPokerList
	return
}

// 转换成一下出牌类型,返回还剩什么手牌
func convertPushPokerList(handPokerList PokerList, findPokers string) (lastHandsPokerList, abovePokerList PokerList) {
	// 记录一下哪些被选了
	var skipMap = make(map[uint32]bool)

	// 一个一个找需要出的牌
	for _, poker := range findPokers {
	out:
		for _, p := range handPokerList {
			// 如果相等 且没被选过
			if string(poker) == string(p.Sn) && !skipMap[p.Index] {
				abovePokerList = append(abovePokerList, p)

				// 记录一下被选了
				skipMap[p.Index] = true

				break out
			}
		}
	}

	// 处理下还剩的手牌
	for _, poker := range handPokerList {
		if !skipMap[poker.Index] {
			lastHandsPokerList = append(lastHandsPokerList, poker)
		}
	}

	// 排个序
	sort.Sort(lastHandsPokerList)
	sort.Sort(abovePokerList)
	return
}

// 处理一下要出的牌和剩下手牌
func diffPushPokerList(handPokerList, findPokers PokerList) (lastHandsPokerList, abovePokerList PokerList) {
	defer func() {
		// 排个序
		sort.Sort(lastHandsPokerList)
		sort.Sort(abovePokerList)
	}()

	// 记录一下哪些被选了
	var skipMap = make(map[uint32]bool)

	// 一个一个找需要出的牌
	for _, poker := range findPokers {
	out:
		for _, p := range handPokerList {
			// 如果相等 且没被选过
			if string(poker.Sn) == string(p.Sn) && !skipMap[p.Index] {
				abovePokerList = append(abovePokerList, p)

				// 记录一下被选了
				skipMap[p.Index] = true

				break out
			}
		}
	}

	// 如果存在牌不存在 那么就表示不出
	if len(skipMap) != len(findPokers) {
		abovePokerList = PokerList{}
		lastHandsPokerList = handPokerList
		return
	}

	// 处理下还剩的手牌
	for _, poker := range handPokerList {
		if !skipMap[poker.Index] {
			lastHandsPokerList = append(lastHandsPokerList, poker)
		}
	}
	return
}
