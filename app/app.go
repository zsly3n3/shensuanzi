package app

import (
	"shensuanzi/commondata"
	"shensuanzi/datastruct/important"
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

func test(r *gin.Engine, handle *handle.AppHandler) {
	r.GET("/app/test", func(c *gin.Context) {
		user_id := commondata.CommonDataInfo.UniqueId()
		app_id := datastruct.SDK_APPID
		c.JSON(200, gin.H{
			"code": handle.AccountGenForIM(user_id, app_id),
		})
	})
}

func RegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	test(r, handle)
}
