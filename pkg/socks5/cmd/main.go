package main

import (
	"github.com/oldbai555/lbtool/pkg/socks5"
	"log"
	"sync"
	"time"
)

// curl -v --proxy socks5://localhost:1080 baidu.com
// curl -v --proxy socks5://username:password@localhost:1080 baidu.com
func main() {

	var mutex sync.Mutex
	server := socks5.Socks5Server{
		Ip:   "localhost",
		Port: 1080,
		Conf: &socks5.Config{
			AuthMethod: socks5.MethodPassWord,
			PasswordChecker: func(username, password string) bool {
				mutex.Lock()
				defer mutex.Unlock()
				log.Printf("username %s,password %s", username, password)
				return true
			},
			TcpTimeout: 5 * time.Second,
		},
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
