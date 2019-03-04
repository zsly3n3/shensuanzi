package app

import (
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

func test(r *gin.Engine, handle *handle.AppHandler) {
	r.GET("/app/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": handle.Test(),
		})
	})
}

func RegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	test(r, handle)
}
