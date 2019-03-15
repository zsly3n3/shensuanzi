package app

import (
	"shensuanzi/commondata"
	"shensuanzi/datastruct"
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

// func checkUserToken(c *gin.Context, handle *handle.AppHandler) (int64, string, bool) {
// 	tokens, isExist := c.Request.Header["Apptoken"]
// 	tf := false
// 	var token string
// 	var userId int64
// 	var isBlackList bool
// 	if isExist {
// 		token = tokens[0]
// 		if token != "" {
// 			userId, tf, isBlackList = handle.IsExistUser(token)
// 			if tf && isBlackList {
// 				// url := eventHandler.GetBlackListRedirect()
// 				url := "http://www.baidu.com"
// 				c.JSON(200, gin.H{
// 					"code": datastruct.Redirect,
// 					"data": url,
// 				})
// 				return -1, "", false
// 			}
// 		}
// 	}
// 	if !tf {
// 		c.JSON(200, gin.H{
// 			"code": datastruct.TokenError,
// 		})
// 	}
// 	return userId, token, tf
// }

func test(r *gin.Engine, handle *handle.AppHandler) {
	url := "/test"
	r.GET(url, func(c *gin.Context) {
		commondata.DeleteOSSFileWithUrl("https://shensuanzi.oss-cn-shenzhen.aliyuncs.com/ft_avatar_dev/110485312978812928.png")
		c.JSON(200, gin.H{
			"code": 0,
		})
	})
}

func isExistUserPhone(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/user/isexistphone/:phone"
	isExistPhone(r, handle, url, false)
}

func getUserDrawCashParams(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/user/drawparams"
	r.GET(url, func(c *gin.Context) {
		_, _, tf := checkFtToken(c, handle)
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

func UserRegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	isExistUserPhone(r, handle)
	getUserDrawCashParams(r, handle)
	test(r, handle)
}
