package main

import (
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/orm"
)

func main() {
	engine, err := orm.NewEngine("mysql", "root:123456@tcp(175.178.156.14:3309)/orm")
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
