package web

import (
	"shensuanzi/commondata"
	"shensuanzi/datastruct"
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

func test(r *gin.Engine) {
	r.GET("/web/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": datastruct.NULLError,
			"data": commondata.CommonDataInfo.UniqueId(),
		})
	})
}

func RegisterRoutes(r *gin.Engine, handle *handle.WebHandler) {
	test(r)
}
