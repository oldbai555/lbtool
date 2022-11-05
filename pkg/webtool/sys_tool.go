package webtool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/extpkg/lbconf/bconf"
	"github.com/oldbai555/lbtool/extpkg/lblog"
	"github.com/oldbai555/lbtool/log"
)

// WebTool 目的 在项目运行中各种中间件都能在此处获取
type WebTool struct {
	ApoC bconf.Config
	Orm  *gorm.DB
	Rdb  *redis.Client

	// Custom 自定义组件
	Log *lblog.Logger
}

// NewWebTool 只支持 apollo
func NewWebTool(conf *ApolloConf, option ...Option) (*WebTool, error) {
	var err error
	lb := &WebTool{}

	// 初始化 apollo 配置中心
	apollo, err := initApollo(conf)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	lb.ApoC = apollo

	// 初始化 mysql redis 等基础组件
	if len(option) == 0 {
		option = []Option{WithGormMysqlOption(), WithRedisOption()}
	}

	// 初始化内置日志服务
	lblog.NewLogger(lblog.SetWriteFile(true))
	lb.Log = lblog.GetLogger()

	// 初始化其他组件
	for _, o := range option {
		err = o.InitConf(apollo)
		if err != nil {
			log.Errorf("err is %v", err)
			return nil, err
		}
		err = o.GenConfTool(lb)
		if err != nil {
			log.Errorf("err is %v", err)
			return nil, err
		}
	}
	return lb, nil
}

func (s *WebTool) GetJson4Apollo(key string, out interface{}) error {
	re, err := s.ApoC.Get(key)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	marshal, err := json.Marshal(re)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	err = json.Unmarshal(marshal, out)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	return nil
}

func (s *WebTool) Infof(ctx context.Context, temp string, args ...interface{}) {
	s.Log.GetCtx(ctx).Info(fmt.Sprintf(temp, args))
}

func (s *WebTool) Warnf(ctx context.Context, temp string, args ...interface{}) {
	s.Log.GetCtx(ctx).Warn(fmt.Sprintf(temp, args))
}

func (s *WebTool) Errorf(ctx context.Context, temp string, args ...interface{}) {
	s.Log.GetCtx(ctx).Error(fmt.Sprintf(temp, args))
}
