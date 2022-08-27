package main

import (
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/utils"
	"github.com/oldbai555/lb/web"
	"net/http"
	"time"
)

var serviceName = "LBW"

func onlyForV1() web.HandlerFunc {
	return func(c *web.Context) error {
		// Start timer
		t := time.Now()
		// Calculate resolution time
		log.Infof("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
		return nil
	}
}

func loadLog() web.HandlerFunc {
	return func(c *web.Context) error {
		log.SetLogHint(c.GetSeq())
		log.SetModuleName(c.GetServerName())
		return nil
	}
}

func init() {
	log.SetEnv(utils.DEV)
}

func main() {
	engine := web.New(serviceName, 12431)
	engine.Use(loadLog())
	engine.GET("/hello", func(c *web.Context) error {
		log.Infof("hello %s", time.Now().Format(utils.DateTimeLayout))
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

	v1 := engine.Group("/v1")
	v1.Use(onlyForV1())
	{
		//v1.GET("/", func(c *web.Context) error {
		//	return c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		//})

		v1.GET("/hello", func(c *web.Context) error {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
			return nil
		})
	}
	v2 := engine.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *web.Context) error {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
			return nil
		})
		v2.POST("/login", func(c *web.Context) error {
			c.JSON(http.StatusOK, map[string]interface{}{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
			return nil
		})

	}

	err := engine.Run()
	if err != nil {
		panic(any(err))
	}
}
