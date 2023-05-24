package utils

import (
	"errors"
	"math/rand"
	"time"
)

/*
 * 产生uint64的不重复数
 * 0xaabbbbbbbbccccdd 其中aa是project id，bbbbbbbb是当前时间戳, cccc是随机数, dd是这一秒内的计数
 * 测试方法:
 * go test -v
 * go test -bench=".*" -parallel 100000
 * cat result*.txt|sort -n |uniq -c | awk '{if($1 != 1){print $0}}' 没有输出说明没有生成重复的id
 * 在我的pc(2.7 GHz Intel Core i5/8 GB 1867 MHz DDR3)上压测结果如下：
 * 10000000	  160 ns/op   5.858s
 * 这种测试条件下，每秒可以产生171w+不重复数字,性能不俗
 */
const (
	ProjectIdBits = 8
	TimestampBits = 32
	RandBits      = 16
	CountBits     = 8
)

var pool = make(chan uint64, 5)
var lastSec uint64 = 0
var lastCount uint8 = 1

func init() {
	go gen()
}

func gen() {
	for {
		current_sec := uint64(time.Now().Unix() & 0x00000000ffffffff)
		if current_sec != lastSec {
			lastCount = 1
			lastSec = current_sec
		}
		c := uint64(lastSec << (RandBits + CountBits))
		rand.Seed(time.Now().UnixNano())
		c += rand.Uint64() & 0x0000000000ffff00
		c += uint64(lastCount)
		lastCount++
		pool <- c
	}
}

func GetUuid4Uint64(projectId uint8) (uint64, error) {
	if c, ok := <-pool; !ok {
		return 0, errors.New("gen uuid fail")
	} else {
		return uint64(projectId)<<(TimestampBits+RandBits+CountBits) + c, nil
	}
}
