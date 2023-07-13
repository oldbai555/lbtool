package sixhecai

import (
	"fmt"
	"github.com/mohae/deepcopy"
	"github.com/oldbai555/lbtool/extpkg/pie/pie"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"strconv"
	"strings"
	"time"
)

var TMList []*TeMa
var SXList = pie.Strings([]string{"兔", "虎", "牛", "鼠", "猪", "狗", "鸡", "猴", "羊", "马", "蛇", "龙"})
var BoSeList = pie.Strings([]string{"绿", "红", "红", "蓝", "蓝", "绿", "绿", "红", "红", "蓝", "绿"})
var DSList = pie.Strings([]string{"双", "单"})
var DXList = pie.Strings([]string{"小", "大"})
var TmStrList = pie.Strings{}
var SXMap = make(map[string][]*TeMa)
var BSMap = make(map[string][]*TeMa)
var DSMap = make(map[string][]*TeMa)
var DXMap = make(map[string][]*TeMa)

func init() {
	sxLen := uint32(len(SXList))
	bsLen := uint32(len(BoSeList))
	dsLen := uint32(len(DSList))
	for i := uint32(0); i < 49; i++ {
		idx := i + 1
		tmStr := fmt.Sprintf("%d", idx)
		if idx < 10 {
			tmStr = fmt.Sprintf("0%d", idx)
		}
		t := &TeMa{
			SX:    SXList[i%sxLen],
			Idx:   idx,
			BS:    BoSeList[idx%bsLen],
			DS:    DSList[idx%dsLen],
			DX:    DXList[idx/25],
			TMStr: tmStr,
		}
		if idx == 10 || idx == 41 {
			t.BS = "蓝"
		}
		TmStrList = append(TmStrList, tmStr)
		TMList = append(TMList, t)
		SXMap[t.SX] = append(SXMap[t.SX], t)
		BSMap[t.BS] = append(BSMap[t.BS], t)
		DSMap[t.DS] = append(DSMap[t.DS], t)
		DXMap[t.DX] = append(DXMap[t.DX], t)
	}

}

type TeMa struct {
	SX    string `json:"sx"`
	Idx   uint32 `json:"idx"`
	BS    string `json:"bs"`
	DS    string `json:"ds"`
	DX    string `json:"dx"`
	TMStr string `json:"tm_str"`
	Je    uint32 `json:"je"`
}

func ShowTeMaList() {
	for _, ma := range TMList {
		log.Infof("tema is %v", ma)
	}
}

func ShowSXMap() {
	for sx, tList := range SXMap {
		log.Infof("sx is %v", sx)
		for _, ma := range tList {
			log.Infof("tema is %v", ma)
		}
	}
}
func ShowBSMap() {
	for bs, tList := range BSMap {
		log.Infof("bs is %v", bs)
		for _, ma := range tList {
			log.Infof("tema is %v", ma)
		}
	}
}
func ShowDXMap() {
	for dx, tList := range DXMap {
		log.Infof("dx is %v", dx)
		for _, ma := range tList {
			log.Infof("tema is %v", ma)
		}
	}
}
func ShowDSMap() {
	for ds, tList := range DSMap {
		log.Infof("ds is %v", ds)
		for _, ma := range tList {
			log.Infof("tema is %v", ma)
		}
	}
}

type Xz struct {
	Je uint32
}

type User struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	TmList     []*TeMa   `json:"tm_list"`
	RecordList []*Record `json:"record_list"`
}

type Record struct {
	Str    string `json:"str"`
	T      string `json:"t"`
	All    uint32 `json:"all"`
	TmList []*TeMa
}

func NewTmList() []*TeMa {
	var m []*TeMa
	for _, s := range TMList {
		m = append(m, deepcopy.Copy(s).(*TeMa))
	}
	return m
}

func NewUser(name string) *User {
	return &User{
		Id:     time.Now().UnixMilli(),
		Name:   name,
		TmList: NewTmList(),
	}
}

func (u *User) SaveXz(str string) error {
	lines := strings.Split(str, "\n")
	var all uint32
	var newRow []string
	list := NewTmList()
	for _, line := range lines {
		newRow = append(newRow, line)
		split := strings.Split(line, "各")
		if len(split) != 2 {
			continue
		}
		xz, err := strconv.Atoi(split[1])
		if err != nil {
			log.Errorf("err is %v", err)
			return err
		}
		res0 := strings.Split(split[0], ",")
		for _, s := range res0 {
			for i := range u.TmList {
				if u.TmList[i].SX == s {
					u.TmList[i].Je += uint32(xz)
					all += uint32(xz)
				} else if u.TmList[i].TMStr == s {
					u.TmList[i].Je += uint32(xz)
					all += uint32(xz)
				}
			}
			for i := range list {
				if list[i].SX == s {
					list[i].Je += uint32(xz)
				} else if list[i].TMStr == s {
					list[i].Je += uint32(xz)
				}
			}
		}
	}
	u.RecordList = append(u.RecordList, &Record{
		Str:    u.Name + " " + strings.Join(newRow, ";") + " " + fmt.Sprintf("(总计%d)", all),
		T:      time.Now().Format(utils.DateTimeLayout),
		All:    all,
		TmList: list,
	})
	return nil
}

func (u *User) Show() {
	for _, ma := range u.TmList {
		log.Infof("ma is %+v", ma)
	}
}

func (u *User) ShowRecord() {
	for _, record := range u.RecordList {
		log.Infof("record is %+v", record)
		for _, ma := range record.TmList {
			log.Infof("ma is %+v", ma)
		}
	}
}
