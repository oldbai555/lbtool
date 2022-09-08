package orm

import (
	"github.com/oldbai555/lb/log"
	"testing"
)

type User struct {
	Id   uint64
	Name string `lborm:"primary_key"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	engine, err := NewEngine(DMYSQL, "root:xxxxxx@tcp(xxxxxx:3306)/orm")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	s := engine.NewSession().Model(&User{})

	err = dropTable(s)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	err = createTable(s)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	table, err := doDescTable(s)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	for _, column := range table.columns {
		log.Infof("column %v", column)
	}
	return
}
