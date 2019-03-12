package app

import (
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

func isExistUserPhone(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/user/isexistphone/:phone"
	isExistPhone(r, handle, url, false)
}

func UserRegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	isExistUserPhone(r, handle)
}
