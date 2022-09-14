package utils

import (
	"fmt"
	"strconv"
	"time"
)

var BeiJinTime = time.FixedZone("Beijing Time", int((8 * time.Hour).Seconds()))

var PRCLocation = BeiJinTime

const (
	TimeLayout     = "15:04:05"
	DateLayout     = "2006-01-02"
	DateTimeLayout = "2006-01-02 15:04:05"
	YYmmDDLayout   = "20060102"
)

const (
	Seconds = 1
	Minutes = 60 * Seconds
	Hours   = 60 * Minutes
	Days    = 24 * Hours
)

var WeekdayMap = map[string]int{
	"周一": 1,
	"周二": 2,
	"周三": 3,
	"周四": 4,
	"周五": 5,
	"周六": 6,
	"周日": 0,
}

func init() {
	SetupTimezone()
}

// SetupTimezone 设置 time 包默认时区为北京时间
func SetupTimezone() {
	time.Local = PRCLocation
}

// Day2Second date 20220402
func Day2Second(date uint32) uint32 {
	if date == 0 {
		return 0
	}
	year := date / 10000
	month := (date % 10000) / 100
	day := date % 100
	if year == 0 || month == 0 || day == 0 {
		return 0
	}
	return uint32(time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Now().Location()).Unix())
}

func TimeNow() uint32 {
	return uint32(time.Now().Unix())
}

// 格式化时间 yyyy-MM-dd HH:mm:ss

func DateYYmmDDhhMMSS() string {
	return time.Now().Format("20060102150405")
}

func DateFormat(temp string) string {
	return time.Now().Format(temp)
}

// 时间戳转换为 20220402 的形式
func coverUnix2Date(beginAt, endAt uint32) (uint32, uint32, error) {
	tmBegin := time.Unix(int64(beginAt), 0).Format("20060102")
	tmEndAt := time.Unix(int64(endAt), 0).Format("20060102")
	begin, err := strconv.Atoi(tmBegin)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return 0, 0, err
	}
	end, err := strconv.Atoi(tmEndAt)
	if err != nil {
		fmt.Errorf("err:%v", err)
		return 0, 0, err
	}
	return uint32(begin), uint32(end), nil

}

// GetDateList 获取需要统计的日期
// case params startDate 20220402
// case params endDate 20220405
func GetDateList(startDate, endDate uint32) (out []uint32) {
	if startDate > endDate {
		return out
	}
	st := time.Unix(int64(Day2Second(startDate)), 0)
	startDate = uint32(0)
	for startDate != endDate {
		startDate64, _ := strconv.ParseUint(st.Format("20060102"), 10, 64)
		startDate = uint32(startDate64)
		out = append(out, startDate)
		st = st.Add(time.Second * 3600 * 24)
		if startDate >= endDate {
			return
		}
	}
	return
}
