package app

import (
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

/*
func checkFtPhone(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/checkphone/:phone"
	r.GET(url, func(c *gin.Context) {
		phone := c.Param("phone")
		if phone == "" {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := handle.checkPhone(phone, true)
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
*/

func RegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	//checkFtPhone(r, handle)
}
