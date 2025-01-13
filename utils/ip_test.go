package utils

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestExternalIP(t *testing.T) {
	ip, err := ExternalIP()
	if err != nil {
		fmt.Printf("err is : %v", err)
		return
	}
	fmt.Println(ip)
}

func Test_getIpFromAddr(t *testing.T) {
	type args struct {
		addr net.Addr
	}
	tests := []struct {
		name string
		args args
		want net.IP
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIpFromAddr(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIpFromAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
