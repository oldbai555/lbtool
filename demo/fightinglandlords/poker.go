package fightinglandlords

import (
	"math/rand"
	"sort"
	"time"
)

// Poker 扑克牌
type Poker struct {
	Index uint32  // 牌的序号
	Type  Decors  // 花色
	Sn    PokerSn // 牌号
}

// NewPokers 初始化一副斗地主的扑克牌
func NewPokers() []*Poker {
	var pokerList, shuffleList []*Poker

	// 初始化牌
	for i := uint32(0); i < pokerCount; i++ {
		pokerList = append(pokerList, &Poker{
			Index: i,
			Type:  decorsMap[i],
			Sn:    PokerSn(pokerInt2StrMap[i]),
		})
	}

	// 打乱顺序
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(pokerList) > 0 {
		n := len(pokerList)
		randIndex := r.Intn(n)
		pokerList[n-1], pokerList[randIndex] = pokerList[randIndex], pokerList[n-1]
		shuffleList = append(shuffleList, pokerList[n-1])
		pokerList = pokerList[:n-1]
	}

	return shuffleList
}

// ToPokers 转换为扑克牌的牌号 indexList 牌的索引
func ToPokers(indexList ...uint32) string {
	res := make([]byte, 0)
	for _, poker := range indexList {
		if poker == smallKing {
			res = append(res, PokerSnMinKing) // 小王
		} else if poker == bigKing {
			res = append(res, PokerSnMaxKing) // 大王
		} else {
			res = append(res, totalCards[(poker/4)%totalNumberPokers]) // 正常序号牌
		}
	}
	return string(res)
}

// ToDecors 转换花色
func ToDecors(index uint32) Decors {
	if index >= smallKing {
		return DecorsNil
	}

	// 处理花色
	switch Decors(index % 4) {
	case DecorsSpade:
		return DecorsSpade
	case DecorsRedHeart:
		return DecorsRedHeart
	case DecorsPlumBlossom:
		return DecorsPlumBlossom
	case DecorsDiamond:
		return DecorsDiamond
	}
	return DecorsNil
}

// SortStr 排序一下牌型
func SortStr(pokers string) string {
	runeArr := make([]int, 0)
	for _, s := range pokers {
		runeArr = append(runeArr, int(s))
	}
	sort.Ints(runeArr)
	res := make([]byte, 0)
	for _, v := range runeArr {
		res = append(res, byte(v))
	}
	return string(res)
}
