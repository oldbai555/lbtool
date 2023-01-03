package webtool

import (
	"fmt"
	mysql "github.com/oldbai555/driver-mysql"
	"github.com/oldbai555/gorm"
	ormlog "github.com/oldbai555/gorm/logger"
	"github.com/oldbai555/lbtool/extpkg/lbconf/bconf"
	"github.com/oldbai555/lbtool/extpkg/lblog"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	syslog "log"
	"time"
)

const defaultApolloMysqlPrefix = "mysql"
const defaultDatabase = "biz"

type GormMysqlConf struct {
	Addr        string `json:"addr"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	TablePrefix string `json:"table_prefix"`
}

func (m *GormMysqlConf) InitConf(apollo bconf.Config) error {
	var v GormMysqlConf
	err := getJson4Apollo(apollo, defaultApolloMysqlPrefix, &v)
	if err != nil {
		log.Errorf(fmt.Sprintf("err is : %v", err))
		return err
	}
	log.Infof("init mysql successfully")
	m.Password = v.Password
	m.Port = v.Port
	m.Addr = v.Addr
	m.Username = v.Username
	if v.TablePrefix == "" {
		m.TablePrefix = "lb_"
	}
	return err
}

func (m *GormMysqlConf) GenConfTool(tool *WebTool, modelObj ...interface{}) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.Username, m.Password, m.Addr, m.Port, defaultDatabase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: gorm.NamingStrategy{
			TablePrefix:   m.TablePrefix, // 指点表名前缀
			SingularTable: true,          // 是否单表，命名是否复数
			NoLowerCase:   false,         // 是否关闭驼峰命名
		},
		NowFunc: func() int32 {
			return int32(time.Now().Unix())
		},
		PrepareStmt: true, // 预编译 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
		Logger:      ormlog.Default.LogMode(ormlog.Info),
	})
	if err != nil {
		log.Errorf(fmt.Sprintf("err is : %v", err))
		return err
	}

	// 自动迁移表 指定建表语句的尾缀
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin").AutoMigrate(modelObj...)
	if err != nil {
		log.Errorf(fmt.Sprintf("err is : %v", err))
		return err
	}

	//  日志配制
	ormlog.New(
		syslog.New(
			lblog.NewFileWriteAsyncer(fmt.Sprintf("go-lb/logs/%s-%s.Log", "mysql", utils.DateFormat(utils.YYmmDDLayout))),
			"\r\n",
			syslog.LstdFlags,
		),
		ormlog.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  ormlog.Silent, // 日志级别
			IgnoreRecordNotFoundError: false,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		})

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Infof("init Orm engine successfully")
	tool.Orm = db
	return nil
}
