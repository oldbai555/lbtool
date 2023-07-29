package event

import (
	"fmt"
	"testing"
)

func TestReg(t *testing.T) {
	Reg(1, func(msg IMsg) error {
		fmt.Println("1-1", msg.GetValue())
		return nil
	})
	Reg(1, func(msg IMsg) error {
		fmt.Println("1-2", msg.GetValue())
		return nil
	})
	Reg(2, func(msg IMsg) error {
		fmt.Println("2-1", msg.GetValue())
		return nil
	})
	Trigger(1, NewMsg("111"))
	Trigger(2, NewMsg("222"))
}

func TestTrigger(t *testing.T) {
	type args struct {
		t Type
		m IMsg
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Trigger(tt.args.t, tt.args.m)
		})
	}
}
