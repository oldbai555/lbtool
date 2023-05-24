package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Map2Struct 用map填充结构
func Map2Struct(data map[string]interface{}, obj interface{}) error {
	for k, v := range data {
		err := setField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// setField 用map的值替换结构的值
func setField(obj interface{}, name string, value interface{}) error {
	structValue := ValueOfData(obj)                   //结构体属性值
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return fmt.Errorf(fmt.Sprintf("No such field: %s in obj", name))
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf(fmt.Sprintf("Cannot set %s field value", name))
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := ValueOfData(value)                  //map值的反射值

	var err error
	if structFieldType != val.Type() {
		val, err = typeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

// typeConversion 类型转换
func typeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 8)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 32)
		return reflect.ValueOf(int32(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 32)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "uint32" {
		i, err := strconv.ParseInt(value, 10, 32)
		return reflect.ValueOf(uint32(i)), err
	} else if ntype == "uint64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(uint64(i)), err
	}

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
