package db

import (
	"fmt"
	"shensuanzi/datastruct"
	"shensuanzi/log"
	"time"

	"github.com/go-xorm/xorm"
)

func (handle *DBHandler) GetDirectDownloadApp() string {
	engine := handle.mysqlEngine
	column_name := "down_load_u_i_addr"
	sql := "select down_load_u_i_addr from domain_info where id = ?"
	results, err := engine.Query(sql, datastruct.DefaultId)
	if err != nil || len(results) <= 0 {
		return ""
	}
	return string(results[0][column_name][:])
}

func (handle *DBHandler) IsExistPhone(phone string, isFT bool) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	var sql string
	if isFT {
		sql = "select 1 from cold_f_t_info where phone = ? limit 1"
	} else {
		sql = "select 1 from cold_user_info where phone = ? limit 1"
	}
	results, err := engine.Query(sql, phone)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	count := len(results)
	isExist := false
	if count > 0 {
		isExist = true
	}
	return isExist, datastruct.NULLError
}

func (handle *DBHandler) IsExistNickName(nickname string) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select 1 from cold_f_t_info where nick_name = ? limit 1"
	results, err := engine.Query(sql, nickname)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	count := len(results)
	isExist := false
	if count > 0 {
		isExist = true
	}
	return isExist, datastruct.NULLError
}
func (handle *DBHandler) GetFtMarkInfo() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select ftm.desc from f_t_mark_info ftm order by id desc"
	results, err := engine.Query(sql)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	arr := make([]string, 0, len(results))
	for _, v := range results {
		arr = append(arr, string(v["desc"][:]))
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) FtRegister(body *datastruct.FTRegisterBody) datastruct.CodeType {
	return ftRegister(handle.mysqlEngine, body, nil)
}

func (handle *DBHandler) FtRegisterWithID(body *datastruct.FTRegisterWithIDBody) datastruct.CodeType {
	return ftRegister(handle.mysqlEngine, &body.FTRegisterBody, body)
}

func ftRegister(engine *xorm.Engine, body *datastruct.FTRegisterBody, IDbody *datastruct.FTRegisterWithIDBody) datastruct.CodeType {
	nowTime := time.Now().Unix()
	cold_ft := new(datastruct.ColdFTInfo)
	cold_ft.CreatedAt = nowTime
	cold_ft.Introduction = body.Desc
	cold_ft.NickName = body.NickName
	cold_ft.Phone = body.Phone
	cold_ft.Pwd = body.Pwd
	cold_ft.Avatar = body.Avatar
	cold_ft.Registration = body.Platform
	cold_ft.AuthState = datastruct.Authing

	if IDbody != nil {
		cold_ft.IdentityCard = IDbody.Identity
		cold_ft.IdFrontCover = IDbody.IdFrontCover
		cold_ft.IdBehindCover = IDbody.IdBehindCover
		cold_ft.ActualName = IDbody.ActualName
	}

	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	_, err := session.InsertOne(cold_ft)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtRegister insert ColdFTInfo :%s", err.Error())
		rollbackError(str, session)
		return datastruct.NickNamePhoneIsUsed
	}

	hot_ft := new(datastruct.HotFTInfo)
	hot_ft.Id = cold_ft.Id
	hot_ft.Mark = body.Mark

	auth := new(datastruct.Authentication)
	auth.FTId = cold_ft.Id
	auth.IsIdCard = false
	auth.IsCP = false
	auth.IsHR = false

	_, err = session.Insert(hot_ft, auth)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtRegister insert HotFTInfo and Authentication :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtRegister Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func rollbackError(err_str string, session *xorm.Session) {
	log.Error("will rollback,err_str:%v", err_str)
	session.Rollback()
}
