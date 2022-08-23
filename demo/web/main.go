package main

import (
	"github.com/oldbai555/lb/comm"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/web"
	"net/http"
	"time"
)

func main() {
	engine := web.New()
	engine.GET("/hello", func(c *web.Context) error {
		log.Infof("hello %s", time.Now().Format(comm.DateTimeLayout))
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		return nil
	})

	engine.GET("/hello/:name", func(c *web.Context) error {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		return nil
	})

	engine.GET("/assets/*filepath/no", func(c *web.Context) error {
		c.JSON(http.StatusOK, web.H{"filepath": c.Param("filepath")})
		return nil
	})

	err := engine.Run(12431)
	if err != nil {
		panic(any(err))
	}
}
