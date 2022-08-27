package main

import (
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/utils"
	"github.com/oldbai555/lb/web"
	"net/http"
	"time"
)

func init() {
	log.SetModuleName("LBW")
}

func main() {
	engine := web.New("myWebTest", 12431)
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
	{
		v1.GET("/", func(c *web.Context) error {
			return c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

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
