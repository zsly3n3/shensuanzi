package handle

import (
	"shensuanzi/commondata"
	"shensuanzi/datastruct"

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
	if err != nil || body.Phone == "" || body.Pwd == "" || body.NickName == "" || body.Avatar == "" || body.Mark == "" || len(body.Desc) < 5 {
		return datastruct.ParamError
	}
	return app.dbHandler.FtRegister(&body)
}

func (app *AppHandler) FtRegisterWithID(c *gin.Context) datastruct.CodeType {
	var body datastruct.FTRegisterWithIDBody
	err := c.BindJSON(&body)
	if err != nil || body.Phone == "" || body.Pwd == "" || body.NickName == "" || body.Avatar == "" || body.Mark == "" || len(body.Desc) < 5 || body.ActualName == "" || body.Identity == "" ||body.IdFrontCover == "" || body.IdBehindCover == "" {
		return datastruct.ParamError
	}
	return app.dbHandler.FtRegisterWithID(&body)
}
