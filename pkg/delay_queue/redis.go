package delay_queue

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rdb *redis.Client

// NewRedisClient 新建 redis client 连接
func NewRedisClient(addr, password string, db int, timeout time.Duration) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password, // no password set
		DB:          db,       // use default DB
		ReadTimeout: timeout,
	})
}

// 执行redis命令, 执行完成后连接自动放回连接池
func execRedisCommand(command string, args ...interface{}) (interface{}, error) {
	err := Rdb.Do(context.TODO(), command, args).Err()
	return nil, err
}
