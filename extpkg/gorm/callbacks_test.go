package gorm_test

import (
	"github.com/oldbai555/lbtool/extpkg/gorm"
	"reflect"
	"sync"
	"testing"
)

type UserWithCallback struct{}

func (UserWithCallback) BeforeSave(*DB) error {
	return nil
}

func (UserWithCallback) AfterCreate(*DB) error {
	return nil
}

func TestCallback(t *testing.T) {
	user, err := gorm.Parse(&UserWithCallback{}, &sync.Map{}, gorm.NamingStrategy{})
	if err != nil {
		t.Fatalf("failed to parse user with callback, got error %v", err)
	}

	for _, str := range []string{"BeforeSave", "AfterCreate"} {
		if !reflect.Indirect(reflect.ValueOf(user)).FieldByName(str).Interface().(bool) {
			t.Errorf("%v should be true", str)
		}
	}

	for _, str := range []string{"BeforeCreate", "BeforeUpdate", "AfterUpdate", "AfterSave", "BeforeDelete", "AfterDelete", "AfterFind"} {
		if reflect.Indirect(reflect.ValueOf(user)).FieldByName(str).Interface().(bool) {
			t.Errorf("%v should be false", str)
		}
	}
}
