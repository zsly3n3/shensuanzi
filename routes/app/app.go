package app

import (
	"shensuanzi/datastruct"
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

func isExistFTPhone(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/isexistphone/:phone"
	isExistPhone(r, handle, url, true)
}
func isExistUserPhone(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/user/isexistphone/:phone"
	isExistPhone(r, handle, url, false)
}
func isExistPhone(r *gin.Engine, handle *handle.AppHandler, url string, isFT bool) {
	r.GET(url, func(c *gin.Context) {
		phone := c.Param("phone")
		if phone == "" {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := handle.IsExistPhone(phone, isFT)
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

func isExistFtNickName(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/isexistnickname/:nickname"
	r.GET(url, func(c *gin.Context) {
		nickname := c.Param("nickname")
		if nickname == "" {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := handle.IsExistNickName(nickname)
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
	isExistFTPhone(r, handle)
	isExistUserPhone(r, handle)
	isExistFtNickName(r, handle)
}
