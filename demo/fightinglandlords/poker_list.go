package fightinglandlords

import (
	"github.com/oldbai555/lbtool/utils"
	"strings"
)

type PokerList []*Poker

// Len 重写 Len() 方法
func (a PokerList) Len() int {
	return len(a)
}

// Swap 重写 Swap() 方法
func (a PokerList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less 重写 Less() 方法， 从大到小排序
func (a PokerList) Less(i, j int) bool {
	return a[j].Sn > a[i].Sn // 从小到大
	// return a[j].Index < a[i].Index // 从大到小
}

func (a PokerList) String() string {
	return strings.Join(utils.PluckStringList(a, "Sn"), "")
}

func (a PokerList) IndexList() []uint32 {
	return utils.PluckUint32List(a, "Index")
}
