package main

import (
	"fmt"
	"github.com/oldbai555/lbtool/demo/proto_emicklei"
	"strconv"
	"time"
)

func main() {
	// 2.验证timer只能响应1次
	// timer3 := time.NewTimer(2 * time.Second)
	//
	// fmt.Printf("2秒到,%v", <-timer3.C)

	s := "23.66"
	float, _ := strconv.ParseFloat(s, 10)
	fmt.Println(float)
	format := time.UnixMilli(int64(1665547918542)).Format("2006-01-02T15:04:05.000+0000")
	fmt.Println(format)

	proto_emicklei.NewDefault().Execute()
}

// GetLastXStr 获取最后几个字符
// prefixStr 剩下的字符
// suffixStr 最后几个字符
func GetLastXStr(str string, lastLen int) (prefixStr string, suffixStr string) {
	rs := []rune(str)
	return string(rs[:len(rs)-lastLen]), string(rs[len(rs)-lastLen:])
}

func Print() {
	//  _ooOoo_
	//                               o8888888o
	//                               88" . "88
	//                               (| -_- |)
	//                                O\ = /O
	//                            ____/`---'\____
	//                          .   ' \\| |// `.
	//                           / \\||| 1 |||// \
	//                         / _||||| -9- |||||- \
	//                           | | \\\ 9 /// | |
	//                         | \_| ''\-8-/'' | |
	//                          \ .-\__ `0` ___/-. /
	//                       ___`. .' /--2--\ `. . __
	//                    ."" '< `.___\_<0>_/___.' >'"".
	//                   | | : `- \`.;`\ 2 /`;.`/ - ` : | |
	//                     \ \ `-. \_ __\ /__ _/ .-` / /
	//             ======`-.____`-.___\_____/___.-`____.-'======
	//                                `=---='
	//
	//             .............................................
	//                      佛祖保佑                  永无BUG
	//              佛曰:
	//                      写字楼里写字间，写字间里程序员；
	//                      程序人员写程序，又拿程序换酒钱。
	//                      酒醒只在网上坐，酒醉还来网下眠；
	//                      酒醉酒醒日复日，网上网下年复年。
	//                      但愿老死电脑间，不愿鞠躬老板前；
	//                      奔驰宝马贵者趣，公交自行程序员。
	//                      别人笑我忒疯癫，我笑自己命太贱；
	//                      不见满街漂亮妹，哪个归得程序员？
	//
}
