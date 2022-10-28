package main

import (
	"github.com/oldbai555/lbtool/discard/orm"
	"github.com/oldbai555/lbtool/discard/orm/dialect"
	"github.com/oldbai555/lbtool/log"
)

func main() {
	engine, err := orm.NewEngine(dialect.DMYSQL, "root:xxxxxx@tcp(xxxxxx:3306)/orm")
	defer func() {
		engine.Close()
	}()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	exec, err := engine.NewSession().Raw("INSERT INTO orm_user(`name`) values (?), (?)", "Tom", "Sam").Exec()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Infof("affected is %d", affected)
}
