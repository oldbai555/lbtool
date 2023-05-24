package main

import (
	"io"
	"log"
	"net/http"
	"testing"
)

func TestAntPool_test(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_helloHandler(t *testing.T) {
	// proxyServer监听端口2002，提供HTTP request “/hello” 的转发服务。
	http.HandleFunc("/true/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":2002", nil))
}

func trueHelloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func Test_trueServer(t *testing.T) {
	http.HandleFunc("/hello", trueHelloHandler)
	log.Fatal(http.ListenAndServe(":2003", nil))
}
