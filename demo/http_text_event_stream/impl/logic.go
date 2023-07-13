package impl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/log"
	"net/http"
	"strings"
)

var ifChannelsMapInit = false

var channelsMap = map[string]chan string{}

func initChannelsMap() {
	channelsMap = make(map[string]chan string)
}

func AddChannel(uniKey string) {
	if !ifChannelsMapInit {
		initChannelsMap()
		ifChannelsMapInit = true
	}
	var newChannel = make(chan string)
	channelsMap[uniKey] = newChannel
	log.Infof("Build SSE connection for uni key = " + uniKey)
}

func BuildNotificationChannel(uniKey string, c *gin.Context) {
	var err error

	AddChannel(uniKey)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	w := c.Writer
	flusher, _ := w.(http.Flusher)
	closeNotify, cancel := context.WithCancel(c.Request.Context())
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	go func() {
		<-closeNotify.Done()
		delete(channelsMap, uniKey)
		log.Infof("SSE close for uni key = " + uniKey)
		return
	}()

	_, err = w.Write([]byte(genMsg("--ping--")))
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	flusher.Flush()

	go func() {
		SendNotification(uniKey)
		cancel()
	}()

	for msg := range channelsMap[uniKey] {
		_, err = w.Write([]byte(msg))
		if err != nil {
			log.Errorf("err is %v", err)
			continue
		}
		flusher.Flush()
	}

}

func SendNotification(uniKey string) {
	log.Infof("Send notification to uni key = " + uniKey)
	var msg = "hello world"
	for key := range channelsMap {
		if strings.Contains(key, uniKey) {
			channel := channelsMap[key]
			for i := 0; i < 10; i++ {
				channel <- genMsg(msg)
			}
		}
	}
}

func genMsg(str string) string {
	return fmt.Sprintf("data: %s\n\n", str)
}
