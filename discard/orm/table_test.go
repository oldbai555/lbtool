package orm

import (
	"github.com/oldbai555/lbtool/discard/orm/dialect"
	"github.com/oldbai555/lbtool/log"
	"testing"
)

type User struct {
	Id   uint64 `lborm:"primary_key"`
	Name string
	Age  int
}

type Cart struct {
	Id   uint64 `lborm:"primary_key"`
	Name string
}

func TestSession_CreateTable(t *testing.T) {
	engine, err := NewEngine(dialect.DMYSQL, "root:xxxxxx@tcp(xxxxxx:3306)/orm")
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

	for _, column := range table.Columns {
		log.Infof("column %v", column)
	}
	return
}

func Test_createOrUpdateTable(t *testing.T) {
	engine, err := NewEngine(dialect.DMYSQL, "root:xxxxxx@tcp(xxxxxx:3306)/orm")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	err = createOrUpdateTable(engine.NewSession().Model(&User{}))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
