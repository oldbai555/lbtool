package main

import (
	"lb/comm"
	"lb/example/conf"
	"lb/log"
	"lb/web"
	"net/http"
	"time"
)

func init() {
	if err := conf.SetupSetting(); err != nil {
		log.Errorf("init.setupSetting err: %v", err)
	}
	conf.ValidateConfig(conf.Settings)
	log.SetUpLogger(conf.Settings.Server.Env)
}

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

	engine.GET("/assets/*filepath", func(c *web.Context) error {
		c.JSON(http.StatusOK, web.H{"filepath": c.Param("filepath")})
		return nil
	})

	err := engine.Run(conf.Settings.Server.HttpPort)
	if err != nil {
		panic(any(err))
	}
}
