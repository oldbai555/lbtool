package main

import (
	"github.com/oldbai555/lbtool/demo/http_text_event_stream/impl"
	"github.com/oldbai555/lbtool/log"
)

/*

// 浏览器调用
e = new EventSource('http://127.0.0.1:7891/notification/socket-connection');
e.onmessage = function(event) {
    console.log(event.data);
};

*/

func main() {
	err := impl.Server()
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
}
