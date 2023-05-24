package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lbtool/extpkg/pie/pie"
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

// CacheHelper 只做结构体级别的缓存
type CacheHelper struct {
	redisClient *redis.Client
	exp         time.Duration

	prefix string // 前缀
	corpId uint32 // 额外参数 corpId

	mType      interface{} // 结构体类型
	fieldNames []string    // 只支持基本数据类型的字段 单字段

	isCombination            bool // 是否需要组合字段,不支持全排列，只支持按顺序组合,fieldName按顺序获取用_隔开,例如 ： prefix_[corpId]_[fieldNames...]
	isOnlyUsedCombinationKey bool // 只使用组合字段进行缓存

	isCacheNotCorpIdKey bool // 是否需要缓存不含CorpId的key
}

type NewCacheHelperReq struct {
	RedisClient *redis.Client
	Prefix      string      // 前缀
	MType       interface{} // 结构体类型
	FieldNames  []string    // 只支持基本数据类型的字段生成Key 最多支持三个, 有键值重复的可能

	IsCombination        bool // 是否需要组合字段
	IsOnlyUseCombination bool // 只使用组合字段

	IsCacheHaveCorpIdKey bool // 是否需要缓存所有Key - 指含有corpId 与不含CorpId的 key
}

func NewCacheHelper(req *NewCacheHelperReq) *CacheHelper {
	if req.RedisClient == nil {
		panic("redisGroup nil")
	}

	if req.MType == nil {
		panic("mType nil")
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

	if len(req.FieldNames) > 3 {
		panic("field names too long")
	}

	// 初始化一下前缀
	if req.Prefix == "" {
		req.Prefix = typ.Name()
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

	return &CacheHelper{
		redisClient: req.RedisClient,
		exp:         DefaultCacheExp,

		mType:      req.MType,
		fieldNames: req.FieldNames,
		prefix:     req.Prefix,

		isCombination:            req.IsCombination,
		isOnlyUsedCombinationKey: req.IsOnlyUseCombination,

		isCacheNotCorpIdKey: req.IsCacheHaveCorpIdKey,
	}
}

// WithCorp 设置携带CorpId
func (c *CacheHelper) WithCorp(corpId uint32) *CacheHelper {
	return &CacheHelper{
		redisClient: c.redisClient,
		exp:         c.exp,

		mType:      c.mType,
		fieldNames: c.fieldNames,
		prefix:     c.prefix,
		corpId:     corpId,
	}
}

// SetJson 设置缓存
func (c *CacheHelper) SetJson(ctx context.Context, model interface{}, exp time.Duration) error {
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
	for _, cacheKey := range c.getAllCacheKey(modelValue) {
		err = c.setJson(ctx, cacheKey, model, exp)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}
	}
	log.Infof("====== end set json to cache ======")
	return nil
}

// AsyncSetJson 异步设置单个结构体缓存
func (c *CacheHelper) AsyncSetJson(ctx context.Context, model interface{}, exp time.Duration) {
	routine.Go(ctx, func(ctx context.Context) error {
		err := c.SetJson(ctx, model, exp)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})
}

// GetJson 获取缓存
// params fieldValue 字段值
func (c *CacheHelper) GetJson(ctx context.Context, fieldValue interface{}, opt interface{}) error {
	err := c.getJson(ctx, c.genCacheKey(fieldValue), opt)
	if err != nil {
		log.Errorf("err:%v", err)
	}
	return err
}

// GetJsonByCustomizeKey 获取缓存 通过组合自定义的valueList 按顺序拼接
func (c *CacheHelper) GetJsonByCustomizeKey(ctx context.Context, opt interface{}, valueList ...string) error {
	err := c.getJson(ctx, c.genCacheKey(strings.Join(valueList, "_")), opt)
	if err != nil {
		log.Errorf("err:%v", err)
	}
	return err
}

// DelJson 清缓存
func (c *CacheHelper) DelJson(ctx context.Context, model interface{}) error {
	modelValue, err := c.getTargetModelValue(model)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	for _, cacheKey := range c.getAllCacheKey(modelValue) {
		err = c.redisClient.Del(ctx, cacheKey).Err()
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}
	}
	return nil
}

// AsyncDelJson 异步清缓存
func (c *CacheHelper) AsyncDelJson(ctx context.Context, model interface{}) {
	routine.Go(ctx, func(ctx context.Context) error {
		err := c.DelJson(ctx, model)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})
}

// 获取所有的缓存 key
func (c *CacheHelper) getAllCacheKey(modelValue reflect.Value) []string {
	var keys pie.Strings
	var allStr []string

	// 遍历结构体生成 key
	for _, fieldName := range c.fieldNames {
		mValue := c.getValueToStr(modelValue, fieldName)
		allStr = append(allStr, mValue)
		// 需要缓存不携带CorpId字段
		if c.isCacheNotCorpIdKey {
			keys = append(keys, c.genCacheNotCorpIdKey(mValue))
		}
		keys = append(keys, c.genCacheKey(mValue))
	}

	// 组合字段 且 只需要组合缓存字段
	if c.isCombination && c.isOnlyUsedCombinationKey {
		var res pie.Strings
		// 需要缓存不携带CorpId字段
		if c.isCacheNotCorpIdKey {
			res = append(res, c.genCacheNotCorpIdKey(strings.Join(allStr, "_")))
		}
		res = append(res, c.genCacheKey(strings.Join(allStr, "_")))
		return res.Unique()
	}

	// 组合字段 且 需要缓存携带与不携带CorpId字段
	if c.isCombination && c.isCacheNotCorpIdKey {
		keys = append(keys, c.genCacheNotCorpIdKey(strings.Join(allStr, "_")))
	}

	// 组合字段
	if c.isCombination {
		keys = append(keys, c.genCacheKey(strings.Join(allStr, "_")))
	}

	return keys.Unique()
}

// 获取缓存的key
func (c *CacheHelper) genCacheKey(val interface{}) string {
	if c.corpId > 0 {
		return fmt.Sprintf("cache_helper_%s_%d_%v", c.prefix, c.corpId, val)
	}
	return fmt.Sprintf("cache_helper_%s_%v", c.prefix, val)

}

func (c *CacheHelper) genCacheNotCorpIdKey(val interface{}) string {
	return fmt.Sprintf("cache_helper_%s_%v", c.prefix, val)
}

// 获取对象的值 转换成字符串
func (c *CacheHelper) getValueToStr(modelValue reflect.Value, fieldName string) string {
	field := modelValue.FieldByName(fieldName)
	if reflect.DeepEqual(field, NullValue) {
		panic(fmt.Sprintf("field %s not found", fieldName))
	}
	// 只支持基本数据类型 作为key
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", field.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", field.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", field.Float())
	case reflect.String:
		return fmt.Sprintf("%s", field.String())
	default:
		panic(fmt.Sprintf("field %v unsupported", field))
	}
}

// 获取目标的Value
func (c *CacheHelper) getTargetModelValue(model interface{}) (reflect.Value, error) {
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

// 检查对象类型是否与 cache_helper 指定的 mType 一致
func (c *CacheHelper) checkValType(model interface{}) bool {
	typeOfMType := reflect.TypeOf(c.mType)
	typeOfModel := reflect.TypeOf(model)
	// 首先它俩得都是指针
	if typeOfMType.Kind() != typeOfModel.Kind() {
		return false
	}

	return strings.TrimPrefix("*", reflect.TypeOf(c.mType).String()) == strings.TrimPrefix("*", reflect.TypeOf(model).String())
}

func (c *CacheHelper) setJson(ctx context.Context, key string, j interface{}, exp time.Duration) error {
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

func (c *CacheHelper) getJson(ctx context.Context, key string, j interface{}) error {
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
