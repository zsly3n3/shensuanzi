package app

import (
	"shensuanzi/datastruct"
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

func checkUserToken(c *gin.Context, handle *handle.AppHandler) (int64, string, bool) {
	tokens, isExist := c.Request.Header["Apptoken"]
	tf := false
	var token string
	var userId int64
	var isBlackList bool
	if isExist {
		token = tokens[0]
		if token != "" {
			userId, tf, isBlackList = handle.IsExistUser(token)
			if tf && isBlackList {
				// url := eventHandler.GetBlackListRedirect()
				url := "http://www.baidu.com"
				c.JSON(200, gin.H{
					"code": datastruct.Redirect,
					"data": url,
				})
				return -1, "", false
			}
		}
	}
	if !tf {
		c.JSON(200, gin.H{
			"code": datastruct.TokenError,
		})
	}
	return userId, token, tf
}

func isExistUserPhone(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/user/isexistphone/:phone"
	isExistPhone(r, handle, url, false)
}

func getUserDrawCashParams(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/user/drawparams"
	r.GET(url, func(c *gin.Context) {
		_, _, tf := checkUserToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetDrawCashParams(datastruct.User)
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

// func userRegister(r *gin.Engine, handle *handle.AppHandler) {
// 	url := "/app/ft/register"
// 	r.POST(url, func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"code": handle.userRegister(c),
// 		})
// 	})
// }

// func userRegisterWithDetail(r *gin.Engine, handle *handle.AppHandler) {
// 	url := "/app/ft/registerwithdetail"
// 	r.POST(url, func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"code": handle.userRegisterWithDetail(c),
// 		})
// 	})
// }

func UserRegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	isExistUserPhone(r, handle)
	getUserDrawCashParams(r, handle)
	// userRegister(r, handle)
	// userRegisterWithDetail(r, handle)
}
