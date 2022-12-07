package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/routine"
	"reflect"
	"strings"
	"time"
)

const (
	DefaultFieldNameById = "Id"
	DefaultCacheExp      = 5 * time.Second
	MTypeNotEqual        = 1201
	ErrSetCache          = 1202
)

var NullValue = reflect.Value{}
var ErrRedisException = errors.New("redis exception")
var ErrJsonUnmarshal = lberr.NewErr(6001, "json unmarshal error")

// cacheHelper 只做结构体级别的缓存
type cacheHelper struct {
	redisClient *redis.Client
	exp         time.Duration

	prefix string // 前缀
	mName  string // 结构体名称

	mType      interface{} // 结构体类型
	fieldNames []string    // 只支持基本数据类型的字段
}

type NewCacheHelperReq struct {
	RedisClient *redis.Client `json:"redis_client"`
	MType       interface{}   `json:"m_type"`      // 结构体类型
	Prefix      string        `json:"prefix"`      // 前缀
	FieldNames  []string      `json:"field_names"` // 只支持基本数据类型的字段
}

func NewCacheHelper(req *NewCacheHelperReq) *cacheHelper {
	if req.RedisClient == nil {
		panic("redisClient nil")
	}

	if req.MType == nil {
		panic("mType nil")
	}

	if req.Prefix == "" {
		panic("prefix is empty")
	}

	// 转换一下
	typ := reflect.TypeOf(req.MType)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// 判断是否是结构体
	if typ.Kind() != reflect.Struct {
		panic("element not struct")
	}

	// 初始化一下字段
	if len(req.FieldNames) == 0 {
		req.FieldNames = []string{DefaultFieldNameById}
	}

	// 校验一下字段
	for _, fieldName := range req.FieldNames {
		// 不支持这些字段
		if fieldName == "DeletedAt" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "CorpId" {
			panic(fmt.Sprintf("field %s unsupported", fieldName))
		}

		// 判断字段是否存在
		field, ok := typ.FieldByName(fieldName)
		if !ok {
			panic(fmt.Sprintf("field %s not found", fieldName))
		}
		// 只支持基本数据类型 作为key
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String:
		default:
			panic(fmt.Sprintf("field %s unsupported", fieldName))
		}
	}

	return &cacheHelper{
		redisClient: req.RedisClient,
		exp:         DefaultCacheExp,

		mType:      req.MType,
		fieldNames: req.FieldNames,
		mName:      typ.Name(),
		prefix:     req.Prefix,
	}
}

// SetJson 设置单个结构体缓存
func (c *cacheHelper) SetJson(ctx context.Context, model interface{}, exp time.Duration) error {
	log.Infof("====== init set json to cache ======")

	modelValue, err := c.getTargetModelValue(model)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if exp == 0 {
		exp = c.exp
	}

	// 根据所需的字段遍历存缓存
	for _, fieldName := range c.fieldNames {
		err = c.setJson(ctx, c.genCacheKey4ModelValue(modelValue, fieldName), model, exp)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}
	log.Infof("====== end set json to cache ======")
	return nil
}

// AsyncSetJson 异步设置单个结构体缓存
func (c *cacheHelper) AsyncSetJson(ctx context.Context, model interface{}, exp time.Duration) {
	routine.Go(ctx, func(ctx context.Context) error {
		err := c.SetJson(ctx, model, exp)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})
}

// BatchSetSingleJson 批量设置单个结构体缓存
func (c *cacheHelper) BatchSetSingleJson(ctx context.Context, modelList []interface{}, exp time.Duration) error {
	for _, val := range modelList {
		err := c.SetJson(ctx, val, exp)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}
	return nil
}

// AsyncBatchSetSingleJson 异步批量设置单个结构体缓存
func (c *cacheHelper) AsyncBatchSetSingleJson(ctx context.Context, modelList []interface{}, exp time.Duration) {
	routine.Go(ctx, func(ctx context.Context) error {
		for _, val := range modelList {
			err := c.SetJson(ctx, val, exp)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}
		return nil
	})
}

// GetJson 获取缓存
// params fieldValue 字段值
func (c *cacheHelper) GetJson(ctx context.Context, fieldValue interface{}, opt interface{}) error {
	// 校验一下类型
	if !c.checkValType(opt) {
		return lberr.NewErr(MTypeNotEqual, "类型不一致")
	}

	err := c.getJson(ctx, c.genCacheKey(fieldValue), opt)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

// GetJsonByDefaultTypeOpt 使用默认类型去创建返回
func (c *cacheHelper) GetJsonByDefaultTypeOpt(ctx context.Context, fieldValue interface{}) (interface{}, error) {
	st := reflect.TypeOf(c.mType)
	for st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	result := reflect.New(st).Interface()

	err := c.getJson(ctx, c.genCacheKey(fieldValue), result)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 校验一下类型
	if !c.checkValType(result) {
		return nil, lberr.NewErr(MTypeNotEqual, "类型不一致")
	}
	return result, nil
}

// DelJson 清缓存
func (c *cacheHelper) DelJson(ctx context.Context, model interface{}) error {
	modelValue, err := c.getTargetModelValue(model)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	for _, fieldName := range c.fieldNames {
		err = c.redisClient.Del(ctx, c.genCacheKey4ModelValue(modelValue, fieldName)).Err()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}
	return nil
}

// AsyncDelJson 异步清缓存
func (c *cacheHelper) AsyncDelJson(ctx context.Context, model interface{}) {
	routine.Go(ctx, func(ctx context.Context) error {
		err := c.DelJson(ctx, model)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})
}

func (c *cacheHelper) genCacheKey(val interface{}) string {
	return fmt.Sprintf("cache_helper_%s_%s_%v", c.mName, c.prefix, val)
}

func (c *cacheHelper) genCacheKey4ModelValue(modelValue reflect.Value, fieldName string) string {
	field := modelValue.FieldByName(fieldName)
	if reflect.DeepEqual(field, NullValue) {
		panic(fmt.Sprintf("field %s not found", fieldName))
	}
	// 只支持基本数据类型 作为key
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("cache_helper_%s_%s_%d", c.mName, c.prefix, field.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("cache_helper_%s_%s_%f", c.mName, c.prefix, field.Float())
	case reflect.String:
		return fmt.Sprintf("cache_helper_%s_%s_%s", c.mName, c.prefix, field.String())
	default:
		panic(fmt.Sprintf("field %v unsupported", field))
	}
}

// getTargetModelType 获取目标的Type
func (c *cacheHelper) getTargetModelType(model interface{}) (reflect.Type, error) {
	// 检查一下类型
	if !c.checkValType(model) {
		return nil, lberr.NewErr(MTypeNotEqual, "类型不一致")
	}

	typ := reflect.TypeOf(model)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// 判断是否是结构体
	if typ.Kind() != reflect.Struct {
		return nil, lberr.NewErr(ErrSetCache, "element not struct")
	}

	return typ, nil
}

// getTargetModelType 获取目标的Value
func (c *cacheHelper) getTargetModelValue(model interface{}) (reflect.Value, error) {
	// 检查一下类型
	if !c.checkValType(model) {
		return reflect.Value{}, lberr.NewErr(MTypeNotEqual, "类型不一致")
	}

	typ := reflect.ValueOf(model)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// 判断是否是结构体
	if typ.Kind() != reflect.Struct {
		return reflect.Value{}, lberr.NewErr(ErrSetCache, "element not struct")
	}

	return typ, nil
}

func (c *cacheHelper) checkValType(model interface{}) bool {
	typeOfMType := reflect.TypeOf(c.mType)
	typeOfModel := reflect.TypeOf(model)
	// 首先它俩得都是指针
	if typeOfMType.Kind() != typeOfModel.Kind() {
		return false
	}

	return strings.TrimPrefix("*", reflect.TypeOf(c.mType).String()) == strings.TrimPrefix("*", reflect.TypeOf(model).String())
}

func (c *cacheHelper) setJson(ctx context.Context, key string, j interface{}, exp time.Duration) error {
	val, err := json.Marshal(j)
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}
	// 空串这里先不考虑
	if len(val) == 0 {
		return errors.New("unsupported empty value")
	}
	return c.redisClient.Set(ctx, key, val, exp).Err()
}

func (c *cacheHelper) getJson(ctx context.Context, key string, j interface{}) error {
	val, err := c.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return redis.Nil
		}
		log.Errorf("err:%s", err)
		return ErrRedisException
	}
	err = json.Unmarshal(val, j)
	if err != nil {
		log.Errorf("err:%s", err)
		return ErrJsonUnmarshal
	}
	return nil
}
