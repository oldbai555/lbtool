package conf

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/oldbai555/comm"
	"github.com/spf13/viper"
	"reflect"
	"sync"
)

//======================== 配制映射结构 ===================================

var Settings *config

type config struct {
	Server serverConfig
}

type serverConfig struct {
	HttpPort uint32 `validate:"required,gt=0"`
	Env      string `validate:"required,oneof=PROD DEV TEST DEMO"`
}

// SetupSetting Setup initialize the configuration instance
func SetupSetting() error {
	var err error
	viper.SetConfigName("config")     // name of config file (without extension)
	viper.AddConfigPath("conf")       // optionally look for config in the working directory
	viper.AddConfigPath("../conf")    // optionally look for config in the working directory
	viper.AddConfigPath("../../conf") // optionally look for config in the working directory
	viper.AddConfigPath("/srv")       // optionally look for config in the working directory
	err = viper.ReadInConfig()        // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		fmt.Errorf("missing config.yaml : %s\n", err.Error())
		return err
	}
	Settings = &config{}
	err = viper.Unmarshal(Settings)
	if err != nil {
		fmt.Errorf("parse config.yaml failed : %s\n", err.Error())
		return err
	}
	viper.WatchConfig()
	return nil
}

//===========================================================

//======================== 校验器 ===================================

type Validator struct {
	Once     sync.Once
	Validate *val.Validate
	Trans    ut.Translator
}

func NewValidator() *Validator {
	return &Validator{}
}

// ValidateConfig 校验配置类
func ValidateConfig(c interface{}) {
	if err := NewValidator().ValidateConfigStruct(c); err != nil {
		panic(any(err.(val.ValidationErrors)))
	}
}

// ValidateConfigStruct 校验配制的结构体
func (v *Validator) ValidateConfigStruct(obj interface{}) error {
	if comm.KindOfData(obj) == reflect.Struct {
		v.LazyInit()
		if err := v.Validate.Struct(obj); err != nil {
			fmt.Errorf("err : %v\n", err)
			return err
		}
	}

	return nil
}

func (v *Validator) LazyInit() {
	v.Once.Do(func() {
		v.Validate = val.New()
		v.Validate.SetTagName("validate")
		// 自定义翻译器(内容自定义翻译校验)
		v.Trans, _ = ut.New(en.New(), zh.New()).GetTranslator("zh")
		zh_translations.RegisterDefaultTranslations(v.Validate, v.Trans)

		//自定义错误内容
		v.Validate.RegisterTranslation("required", v.Trans, func(ut ut.Translator) error {
			return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe val.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})

	})
}

//===========================================================
