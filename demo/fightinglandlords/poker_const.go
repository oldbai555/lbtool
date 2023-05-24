package fightinglandlords

const (
	smallKing = 52 // 小王
	bigKing   = 53 // 大王

	totalNumberPokers = 13 // 除去大小王的序号牌个数
	pokerCount        = 54 // 扑克牌总数

	totalCards = "A234567890JQK" // 牌的序号
)

type Decors uint32 // 花色
const (
	DecorsNil         Decors = iota // 无花色
	DecorsSpade                     // 黑桃
	DecorsRedHeart                  // 红心
	DecorsPlumBlossom               // 梅花
	DecorsDiamond                   // 方块
)

const (
	MaxLandlordPokers = 3    // 地主牌数量
	KingRocket        = "Ww" // 王炸
)

type PokerSn string // 牌号
const (
	PokerSnA       = 'A'
	PokerSn2       = '2'
	PokerSn3       = '3'
	PokerSn4       = '4'
	PokerSn5       = '5'
	PokerSn6       = '6'
	PokerSn7       = '7'
	PokerSn8       = '8'
	PokerSn9       = '9'
	PokerSn0       = '0'
	PokerSnJ       = 'J'
	PokerSnQ       = 'Q'
	PokerSnK       = 'K'
	PokerSnMinKing = 'w'
	PokerSnMaxKing = 'W'
)

// key by index

// 花色Map
var decorsMap = make(map[uint32]Decors, pokerCount)

// 牌组下标对应的字符串
var pokerInt2StrMap = make(map[uint32]string, pokerCount)

// 初始化全局牌号和全局牌花色
func init() {
	for index := uint32(0); index < pokerCount; index++ {
		// 处理花色
		decorsMap[index] = ToDecors(index)
		// 处理牌号
		pokerInt2StrMap[index] = ToPokers(index)
	}
}
