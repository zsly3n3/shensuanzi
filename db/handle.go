package db

import (
	"shensuanzi/datastruct"
)

func (handle *DBHandler) GetDirectDownloadApp() string {
	engine := handle.mysqlEngine
	column_name := "down_load_u_i_addr"
	sql := "select ? from domain_info where id = ?"
	results, err := engine.Query(sql, column_name, datastruct.DefaultId)
	if err != nil || len(results) <= 0 {
		return ""
	}
	return string(results[0][column_name][:])
}
