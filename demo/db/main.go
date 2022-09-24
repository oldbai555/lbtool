package main

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/orm"
	"github.com/oldbai555/lbtool/orm/dialect"
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
