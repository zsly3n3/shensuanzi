package db

import (
	"shensuanzi/datastruct"
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

func (handle *DBHandler) CheckPhone(phone string, isFT bool) (interface{}, datastruct.CodeType) {
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
