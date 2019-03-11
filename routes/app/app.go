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
		data, code := handle.FtLogin(c)
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

func ftInfo(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/info"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtInfo(id)
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

func updateFtInfo(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/info"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.UpdateFtInfo(c, id),
		})
	})
}

func updateFtMark(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/mark"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.UpdateFtMark(c, id),
		})
	})
}

func updateFtIntroduction(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/introduce"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.UpdateFtIntroduction(c, id),
		})
	})
}

func getFtIntroduction(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/introduce"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtIntroduction(id)
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

func getFtAutoReply(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/autoreply"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtAutoReply(id)
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

// func ftIsOnline(r *gin.Engine, handle *handle.AppHandler) {
// 	url := "/app/ft/online"
// 	r.POST(url, func(c *gin.Context) {
// 		id, _, tf := checkFtToken(c, handle)
// 		if !tf {
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"code": handle.FtIsOnline(id),
// 		})
// 	})
// }

func checkFtToken(c *gin.Context, handle *handle.AppHandler) (int, string, bool) {
	tokens, isExist := c.Request.Header["Apptoken"]
	tf := false
	var token string
	var ft_id int
	var isBlackList bool
	if isExist {
		token = tokens[0]
		if token != "" {
			ft_id, tf, isBlackList = handle.IsExistFt(token)
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
	return ft_id, token, tf
}

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

func RegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	isExistFtPhone(r, handle)
	isExistUserPhone(r, handle)
	isExistFtNickName(r, handle)
	getFtMarkInfo(r, handle)
	ftRegister(r, handle)
	ftRegisterWithID(r, handle)
	ftLogin(r, handle)
	ftInfo(r, handle)
	updateFtInfo(r, handle)
	updateFtMark(r, handle)
	updateFtIntroduction(r, handle)
	getFtIntroduction(r, handle)
	getFtAutoReply(r, handle)
	// ftIsOnline(r, handle)
}
