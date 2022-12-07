package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lbtool/log"
	"testing"
	"time"
)

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func TestNewCacheHelper(t *testing.T) {
	log.SetLogHint(fmt.Sprintf("%d", time.Now().Unix()))
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	ctx := context.Background()
	helper := NewCacheHelper(&NewCacheHelperReq{
		Prefix:      "bai",
		RedisClient: rdb,
		MType:       &User{},
		FieldNames:  []string{"Id", "Name"},
	})
	testUser := &User{Id: 1, Name: "bai"}
	err := helper.SetJson(ctx, &testUser, 0)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	var newUser User
	err = helper.GetJson(ctx, 1, &newUser)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Infof("new user: %v", newUser)

	err = helper.DelJson(ctx, &newUser)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	err = helper.GetJson(ctx, 1, &newUser)
	if err != nil {
		log.Errorf("err:%v", err)
		time.Sleep(5 * time.Second)
		return
	}
	log.Infof("new user: %v", newUser)
	time.Sleep(5 * time.Second)
}
