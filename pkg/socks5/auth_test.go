package socks5

import (
	"net"
	"reflect"
	"testing"
)

func TestNewClientAuthMessage(t *testing.T) {
	type args struct {
		conn net.Conn
	}
	tests := []struct {
		name    string
		args    args
		want    *ClientAuthMessage
		wantErr bool
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClientAuthMessage(tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientAuthMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientAuthMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth(t *testing.T) {
	type args struct {
		conn net.Conn
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := auth(tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("auth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
