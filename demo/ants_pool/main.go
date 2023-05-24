package main

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/panjf2000/ants/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func init() {
	log.SetBaseDir(utils.GetCurrentAbPath() + "/log")
}

// AntPool 线程池
func AntPool() {
	// 第一种用法
	pool, err := ants.NewPool(10)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	defer func() {
		rErr := pool.ReleaseTimeout(10 * time.Second)
		if rErr != nil {
			log.Errorf("err:%v", rErr)
			return
		}
	}()

	for i := 0; i < 100; i++ {
		var a = i
		err = pool.Submit(func() {
			log.Infof("%d", a)
			time.Sleep(2 * time.Second)
		})
	}

	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// 第二种用法
	poolWF, err := ants.NewPoolWithFunc(10, func(i interface{}) {
		log.Warnf("%v", i)
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	for i := 0; i < 100; i++ {
		var a = i
		err = poolWF.Invoke(fmt.Sprintf("hello + %d", a))
	}
}

// 将request转发给 http://127.0.0.1:2003
func helloHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	path = strings.TrimPrefix(path, "/")
	split := strings.Split(path, "/")
	if len(split) < 2 {
		fmt.Println("<2")
		return
	}

	// todo 需要重新构造一个 request
	r.URL.Path = "/" + split[1]

	// 代理服务
	proxyServer := "http://127.0.0.1:2003"
	target, err := url.Parse(proxyServer)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// NewSingleHostReverseProxy 返回一个新的 ReverseProxy，
	// 将 URLs 请求路由到 target 的指定的 scheme, host, base path 。
	// 如果 target 的 path 是 ”/base” ，client请求的URL是 “/dir”,
	// 则 target 最后转发的请求就是 /base/dir。
	proxy := httputil.NewSingleHostReverseProxy(target)

	// ReverseProxy 是 HTTP Handler， 接收client的 request，
	// 将其发送给另一个server, 并将server 返回的response转发给client。
	proxy.ServeHTTP(w, r)
}
