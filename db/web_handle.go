package db

import (
	"shensuanzi/datastruct"
	"shensuanzi/log"
)

func (handle *DBHandler) GetServerInfo() (*datastruct.WebServerInfoBody, datastruct.CodeType) {
	engine := handle.mysqlEngine
	serverInfo := new(datastruct.ServerInfo)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(serverInfo)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.WebServerInfoBody)
	resp.IsMaintain = serverInfo.IsMaintain
	resp.Version = serverInfo.Version
	resp.GzhAppid = serverInfo.GzhAppid
	resp.KfptAppid = serverInfo.KfptAppid
	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditServerInfo(body *datastruct.WebServerInfoBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	serverInfo := new(datastruct.ServerInfo)
	serverInfo.GzhAppid = body.GzhAppid
	serverInfo.IsMaintain = body.IsMaintain
	serverInfo.KfptAppid = body.KfptAppid
	serverInfo.Version = body.Version
	_, err := engine.Where("id=?", datastruct.DefaultId).Update(serverInfo)
	if err != nil {
		log.Error("EditServerInfo err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}
