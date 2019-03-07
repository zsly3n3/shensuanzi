package app

import (
	"shensuanzi/datastruct"
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

func isExistFtPhone(r *gin.Engine, handle *handle.AppHandler) {
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

func getFtMarkInfo(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/mark"
	r.GET(url, func(c *gin.Context) {
		data, code := handle.GetFtMarkInfo()
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

func ftRegister(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/register"
	r.POST(url, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": handle.FtRegister(c),
		})
	})
}

func ftRegisterWithID(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/registerwithid"
	r.POST(url, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": handle.FtRegisterWithID(c),
		})
	})
}

func ftLogin(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/login"
	r.POST(url, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": handle.FtLogin(c),
		})
	})
}

func RegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	isExistFtPhone(r, handle)
	isExistUserPhone(r, handle)
	isExistFtNickName(r, handle)
	getFtMarkInfo(r, handle)
	ftRegister(r, handle)
	ftRegisterWithID(r, handle)
	ftLogin(r, handle)
}
