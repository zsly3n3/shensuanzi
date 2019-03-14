package handle

import (
	"shensuanzi/commondata"
	"shensuanzi/datastruct"
	"shensuanzi/tools"

	"github.com/gin-gonic/gin"
)

func (app *AppHandler) GetServerInfoFromMemory() (string, bool) {
	serverInfo := commondata.GetServerInfo()
	serverInfo.RWMutex.RLock()
	defer serverInfo.RWMutex.RUnlock()
	return serverInfo.Version, serverInfo.IsMaintain
}

func (app *AppHandler) GetDirectDownloadApp() string {
	return app.dbHandler.GetDirectDownloadApp()
}

func (app *AppHandler) IsExistPhone(phone string, isFT bool) (interface{}, datastruct.CodeType) {
	return app.dbHandler.IsExistPhone(phone, isFT)
}

func (app *AppHandler) IsExistNickName(nickname string) (interface{}, datastruct.CodeType) {
	return app.dbHandler.IsExistNickName(nickname)
}

func (app *AppHandler) GetFtMarkInfo() (interface{}, datastruct.CodeType) {
	return app.dbHandler.GetFtMarkInfo()
}

func (app *AppHandler) FtRegister(c *gin.Context) datastruct.CodeType {
	var body datastruct.FTRegisterBody
	err := c.BindJSON(&body)
	isRemoveFile := false
	var code datastruct.CodeType
	if err != nil || body.Phone == "" || body.Pwd == "" || body.NickName == "" || body.Avatar == "" || body.Mark == "" || len(body.Desc) < 5 {
		isRemoveFile = true
		code = datastruct.ParamError
	} else {
		code = app.dbHandler.FtRegister(&body)
		if code != datastruct.NULLError {
			isRemoveFile = true
		}
	}
	if isRemoveFile {
		go app.deleteRegisterFile(body.Avatar, "", "")
	}
	return code
}

func (app *AppHandler) deleteRegisterFile(avatar string, IdFrontCover string, IdBehindCover string) {
	if avatar != "" {
		commondata.DeleteOSSFileWithUrl(avatar)
	}
	if IdFrontCover != "" {
		commondata.DeleteOSSFileWithUrl(IdFrontCover)
	}
	if IdBehindCover != "" {
		commondata.DeleteOSSFileWithUrl(IdBehindCover)
	}
}

func (app *AppHandler) FtRegisterWithID(c *gin.Context) datastruct.CodeType {
	var body datastruct.FTRegisterWithIDBody
	err := c.BindJSON(&body)
	isRemoveFile := false
	var code datastruct.CodeType
	if err != nil || body.Phone == "" || body.Pwd == "" || body.NickName == "" || body.Avatar == "" || body.Mark == "" || len(body.Desc) < 5 || body.ActualName == "" || len(body.Identity) != 18 || body.IdFrontCover == "" || body.IdBehindCover == "" {
		isRemoveFile = true
		code = datastruct.ParamError
	} else {
		code = app.dbHandler.FtRegisterWithID(&body)
		if code != datastruct.NULLError {
			isRemoveFile = true
		}
	}
	if isRemoveFile {
		go app.deleteRegisterFile(body.Avatar, body.IdFrontCover, body.IdBehindCover)
	}
	return code
}
func (app *AppHandler) FtLogin(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.FtLoginBody
	err := c.BindJSON(&body)
	if err != nil || body.Phone == "" || body.Pwd == "" {
		return nil, datastruct.ParamError
	}
	rs, code := app.dbHandler.FtLogin(&body)
	if code != datastruct.NULLError {
		return nil, code
	}
	ft_redis := new(datastruct.FtRedisData)
	ft_redis.FtId = rs.FtInfo.Id
	ft_redis.Token = rs.Token
	ft_redis.AccountState = rs.FtInfo.AccountState
	conn := app.cacheHandler.GetConn()
	defer conn.Close()
	app.cacheHandler.SetFtToken(conn, ft_redis)
	app.cacheHandler.AddExpire(conn, ft_redis.Token)
	return rs, code
}

func (app *AppHandler) UpdateFtInfo(c *gin.Context, ft_id int) datastruct.CodeType {
	var body datastruct.UpdateFtInfoBody
	err := c.BindJSON(&body)
	if err != nil || body.Avatar == "" || body.NickName == "" {
		return datastruct.ParamError
	}
	return app.dbHandler.UpdateFtInfo(&body, ft_id)
}

func (app *AppHandler) UpdateFtMark(c *gin.Context, ft_id int) datastruct.CodeType {
	var body datastruct.UpdateFtMarkBody
	err := c.BindJSON(&body)
	if err != nil || body.Mark == "" {
		return datastruct.ParamError
	}
	return app.dbHandler.UpdateFtMark(&body, ft_id)
}

func (app *AppHandler) UpdateFtIntroduction(c *gin.Context, ft_id int) datastruct.CodeType {
	var body datastruct.UpdateFtIntroductionBody
	err := c.BindJSON(&body)
	if err != nil || body.Desc == "" || len(body.Imgs) <= 0 {
		return datastruct.ParamError
	}
	return app.dbHandler.UpdateFtIntroduction(&body, ft_id)
}
func (app *AppHandler) UpdateFtAutoReply(c *gin.Context, ft_id int) datastruct.CodeType {
	var body datastruct.UpdateFtAutoReplyBody
	err := c.BindJSON(&body)
	if err != nil || body.AutoReply == "" || len(body.QuickReply) <= 0 {
		return datastruct.ParamError
	}
	return app.dbHandler.UpdateFtAutoReply(&body, ft_id)
}
func (app *AppHandler) FtSubmitIdentity(c *gin.Context, ft_id int) (interface{}, datastruct.CodeType) {
	var body datastruct.FtIdentity
	err := c.BindJSON(&body)
	if err != nil || body.ActualName == "" || body.IdBehindCover == "" || body.IdFrontCover == "" || len(body.Identity) != 18 {
		return nil, datastruct.ParamError
	}
	return app.dbHandler.FtSubmitIdentity(&body, ft_id)
}

func (app *AppHandler) GetFtIntroduction(ft_id int) (interface{}, datastruct.CodeType) {
	return app.dbHandler.GetFtIntroduction(ft_id)
}

func (app *AppHandler) GetFtAutoReply(ft_id int) (interface{}, datastruct.CodeType) {
	return app.dbHandler.GetFtAutoReply(ft_id)
}

func (app *AppHandler) GetAppraised(c *gin.Context, ft_id int) (interface{}, datastruct.CodeType) {
	pageIndex := tools.StringToInt(c.Param("pageindex"))
	pageSize := tools.StringToInt(c.Param("pagesize"))
	if pageIndex <= 0 || pageSize <= 0 {
		return nil, datastruct.ParamError
	}
	return app.dbHandler.GetAppraised(ft_id, pageIndex, pageSize)
}

func (app *AppHandler) GetFtSystemMsg(c *gin.Context, ft_id int) (interface{}, datastruct.CodeType) {
	pageIndex := tools.StringToInt(c.Param("pageindex"))
	pageSize := tools.StringToInt(c.Param("pagesize"))
	if pageIndex <= 0 || pageSize <= 0 {
		return nil, datastruct.ParamError
	}
	return app.dbHandler.GetFtSystemMsg(ft_id, pageIndex, pageSize)
}

func (app *AppHandler) GetFtDndList(c *gin.Context, ft_id int) (interface{}, datastruct.CodeType) {
	pageIndex := tools.StringToInt(c.Param("pageindex"))
	pageSize := tools.StringToInt(c.Param("pagesize"))
	if pageIndex <= 0 || pageSize <= 0 {
		return nil, datastruct.ParamError
	}
	return app.dbHandler.GetFtDndList(ft_id, pageIndex, pageSize)
}

func (app *AppHandler) RemoveFtDndList(c *gin.Context, ft_id int) datastruct.CodeType {
	var body datastruct.RemoveWithIdBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 {
		return datastruct.ParamError
	}
	return app.dbHandler.RemoveFtDndList(body.Id, ft_id)
}

func (app *AppHandler) EditProduct(c *gin.Context, ft_id int) (interface{}, datastruct.CodeType) {
	var body datastruct.EditProductBody
	err := c.BindJSON(&body)
	if err != nil || body.Price < 20 || body.Price > 3000 || body.ProductDesc != "" || body.ProductName != "" {
		return nil, datastruct.ParamError
	}
	return app.dbHandler.EditProduct(&body, ft_id)
}

func (app *AppHandler) GetProduct(ft_id int) (interface{}, datastruct.CodeType) {
	return app.dbHandler.GetProduct(ft_id)
}

func (app *AppHandler) RemoveProduct(c *gin.Context, ft_id int) datastruct.CodeType {
	var body datastruct.RemoveWithIdBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 {
		return datastruct.ParamError
	}
	return app.dbHandler.RemoveProduct(body.Id, ft_id)
}

func (app *AppHandler) SortProducts(c *gin.Context) datastruct.CodeType {
	var body datastruct.RemoveWithIdsBody
	err := c.BindJSON(&body)
	if err != nil || len(body.Ids) <= 0 {
		return datastruct.ParamError
	}
	return app.dbHandler.SortProducts(body.Ids)
}

func (app *AppHandler) CreateFakeAppraised(c *gin.Context, ft_id int) datastruct.CodeType {
	var body datastruct.FakeAppraisedBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 || len(body.Desc) > 50 || body.Desc == "" || body.Score <= 0 || body.Time <= 0 {
		return datastruct.ParamError
	}
	return app.dbHandler.CreateFakeAppraised(body, ft_id)
}

func (app *AppHandler) GetFtUnReadMsgCount(ft_id int) (interface{}, datastruct.CodeType) {
	return app.dbHandler.GetFtUnReadMsgCount(ft_id)
}

func (app *AppHandler) GetFtInfo(ft_id int) (interface{}, datastruct.CodeType) {
	return app.dbHandler.GetFtInfo(ft_id)
}

func (app *AppHandler) GetAllFtOrder(c *gin.Context, ft_id int) (interface{}, datastruct.CodeType) {
	pageIndex := tools.StringToInt(c.Param("pageindex"))
	pageSize := tools.StringToInt(c.Param("pagesize"))
	if pageIndex <= 0 || pageSize <= 0 {
		return nil, datastruct.ParamError
	}
	return app.dbHandler.GetAllFtOrder(ft_id, pageIndex, pageSize)
}

func (app *AppHandler) IsExistFt(token string) (int, bool, bool) {
	conn := app.cacheHandler.GetConn()
	defer conn.Close()
	var ft_id int
	var tf bool
	var isBlackList bool
	ft_id, tf, isBlackList = app.cacheHandler.IsExistFtWithConn(conn, token)
	if !tf {
		var ft_data *datastruct.FtRedisData
		ft_data, tf = app.dbHandler.GetFtDataWithToken(token)
		if tf {
			if ft_data.AccountState == datastruct.BlackList {
				isBlackList = true
			}
			ft_id = ft_data.FtId
			app.cacheHandler.SetFtToken(conn, ft_data)
			app.cacheHandler.AddExpire(conn, token)
		}
	} else {
		app.cacheHandler.AddExpire(conn, token)
	}
	return ft_id, tf, isBlackList
}

func (app *AppHandler) IsExistUser(token string) (int64, bool, bool) {
	conn := app.cacheHandler.GetConn()
	defer conn.Close()
	var userId int64
	var tf bool
	var isBlackList bool
	userId, tf, isBlackList = app.cacheHandler.IsExistUserWithConn(conn, token)
	if !tf {
		// var user *datastruct.UserInfo
		// user, tf = handle.dbHandler.GetUserDataWithToken(token)
		// if tf {
		// 	userId = user.Id
		// 	isBlackList = tools.IntToBool(user.IsBlackList)
		// 	handle.cacheHandler.SetUserAllData(conn, user)
		// 	handle.cacheHandler.AddExpire(conn, token)
		// }
	} else {
		//app.cacheHandler.AddExpire(conn, token)
	}
	return userId, tf, isBlackList
}

// func (app *AppHandler) FtIsOnline(ft_id int) datastruct.CodeType {
// 	conn := app.cacheHandler.GetConn()
// 	defer conn.Close()
// 	return app.cacheHandler.SetFtOnline(conn, ft_id)
// }
