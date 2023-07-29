package worker

import (
	"context"
	"github.com/oldbai555/lbtool/log"
	"reflect"
	"sync"
	"testing"
)

func TestBaseHandlerMgr_Call(t *testing.T) {
	worker := NewWorker(1024, "test")
	worker.Register(1, func(ctx context.Context, i interface{}) error {
		log.Infof("hhhh i %v", i)
		return nil
	})
	worker.Start(context.Background())

	for i := 0; i < 10000; i++ {
		err := worker.Send(1, i)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}
	}
	select {}
}

func TestBaseHandlerMgr_Register(t *testing.T) {
	type fields struct {
		fnM   map[Type]DoHandlerFn
		Mutex sync.Mutex
	}
	type args struct {
		t  Type
		fn DoHandlerFn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &BaseHandlerMgr{
				fnM:   tt.fields.fnM,
				Mutex: tt.fields.Mutex,
			}
			e.Register(tt.args.t, tt.args.fn)
		})
	}
}

func TestBaseMsg_GetType(t *testing.T) {
	type fields struct {
		Typ Type
		Val interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Type
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &BaseMsg{
				Typ: tt.fields.Typ,
				Val: tt.fields.Val,
			}
			if got := i.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMsg_GetValue(t *testing.T) {
	type fields struct {
		Typ Type
		Val interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &BaseMsg{
				Typ: tt.fields.Typ,
				Val: tt.fields.Val,
			}
			if got := i.GetValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
