package app

import (
	"shensuanzi/datastruct"
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

func checkFTPhone(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/checkphone/:phone"
	checkPhone(r, handle, url, true)
}
func checkUserPhone(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/user/checkphone/:phone"
	checkPhone(r, handle, url, false)
}

func checkPhone(r *gin.Engine, handle *handle.AppHandler, url string, isFT bool) {
	r.GET(url, func(c *gin.Context) {
		phone := c.Param("phone")
		if phone == "" {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := handle.CheckPhone(phone, isFT)
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
	checkFTPhone(r, handle)
	checkUserPhone(r, handle)
}
