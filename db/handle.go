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
func (handle *DBHandler) GetFtMarkInfo(nickname string) (interface{}, datastruct.CodeType) {
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
