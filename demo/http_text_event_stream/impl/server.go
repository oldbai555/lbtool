package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/log"
	"net/http"
)

func Server() error {
	engine := gin.New()
	engine.Use(Cors())
	engine.GET("/notification/socket-connection", SocketConnection)
	err := engine.Run("127.0.0.1:7891")
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func SocketConnection(ctx *gin.Context) {
	BuildNotificationChannel("1", ctx)
}

// Cors 跨域配制
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		method := c.Request.Method
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
