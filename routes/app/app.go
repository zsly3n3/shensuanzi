package app

import (
	"shensuanzi/datastruct"
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

func getTest(r *gin.Engine, handle *handle.AppHandler) {
	r.GET("/app/gettest", func(c *gin.Context) {
		data, code := handle.GetTest()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func RegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	test(r, handle)
	getTest(r, handle)
}
