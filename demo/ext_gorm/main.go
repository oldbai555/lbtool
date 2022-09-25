package main

import (
	"fmt"
)

import (
	"github.com/oldbai555/lbtool/extpkg/gorm"
	ormlog "github.com/oldbai555/lbtool/extpkg/gorm/logger"
	"github.com/oldbai555/lbtool/extpkg/gorm_mysql"
	"github.com/oldbai555/lbtool/extpkg/lblog"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	syslog "log"
	"time"
)

func main() {
	db, err := InitOrmEngine()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	db.Create(&User{
		Username: "admin",
		Password: "123456",
	})
	var u User
	err = db.Where("username = ?", "admin").First(&u).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	err = db.Model(&User{}).Where("username = ?", "admin").Update("password", "135246").Error
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	log.Infof("u is %v", u)
	err = db.Where("username = ?", "admin").Delete(&User{}).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Infof("delete")
}

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int32
	UpdatedAt int32
	DeletedAt int32 `gorm:"index;default:0"`

	Username string `json:"username" gorm:"varchar(25); not null; comment('账号')"`
	Password string `json:"password" gorm:"varchar(25); not null; comment('密码')"`
}

// InitOrmEngine https://gorm.io/zh_CN/docs/connecting_to_the_database.html
func InitOrmEngine() (*gorm.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "0", "0", "0", 0, "0")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: gorm.NamingStrategy{
			TablePrefix:   "blog_", // 指点表名前缀
			SingularTable: true,    // 是否单表，命名是否复数
			NoLowerCase:   false,   // 是否关闭驼峰命名
		},
		NowFunc: func() int32 {
			return int32(time.Now().Unix())
		},
		PrepareStmt: true, // 预编译 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
		Logger:      ormlog.Default.LogMode(ormlog.Info),
	})
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	// 自动迁移表 指定建表语句的尾缀
	err = db.
		Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin").
		AutoMigrate(
			&User{},
		)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	//  日志配制
	ormlog.New(
		syslog.New(
			lblog.NewFileWriteAsyncer(fmt.Sprintf("blog_api/logs/%s-%s.log", "mysql", utils.DateFormat(utils.YYmmDDLayout))),
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

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
