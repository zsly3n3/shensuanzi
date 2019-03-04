package db

import (
	"fmt"
	"shensuanzi/conf"
	"shensuanzi/datastruct"
	"shensuanzi/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DBHandler struct {
	mysqlEngine *xorm.Engine
}

func CreateDBHandler() *DBHandler {
	dbHandler := new(DBHandler)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", conf.Server.DB_UserName, conf.Server.DB_Pwd, conf.Server.DB_IP, conf.Server.DB_Name)
	engine, err := xorm.NewEngine("mysql", dsn)
	errhandle(err)
	err = engine.Ping()
	errhandle(err)
	//日志打印SQL
	engine.ShowSQL(true)
	//设置连接池的空闲数大小
	maxIdleConns := 40
	if conf.Common.Mode == conf.Debug {
		maxIdleConns = 5
	}
	engine.SetMaxIdleConns(maxIdleConns)
	SyncDB(engine)
	initData(engine)
	dbHandler.mysqlEngine = engine
	//go timerTask(engine)
	return dbHandler
}

func SyncDB(engine *xorm.Engine) {
	arr := make([]interface{}, 0)
	var err error

	arr = append(arr, new(datastruct.ColdUserInfo))
	arr = append(arr, new(datastruct.HotUserInfo))
	arr = append(arr, new(datastruct.WXUserInfo))
	arr = append(arr, new(datastruct.AdInfo))

	err = engine.Sync2(arr...)
	errhandle(err)

	// role := createRoleData()
	// _, err = engine.Insert(&role)
	// errhandle(err)
	// webUser := createLoginData(datastruct.AdminLevelID)
	// _, err = engine.Insert(&webUser)
	// errhandle(err)
}

/*
func createRoleData() []datastruct.Role {
	admin := datastruct.Role{
		Id:   datastruct.AdminLevelID,
		Desc: "admin",
	}
	guest := datastruct.Role{
		Id:   datastruct.NormalLevelID,
		Desc: "normal",
	}

	return []datastruct.Role{admin, guest}
}

func createLoginData(adminLevelID int) []datastruct.WebUser {
	now_time := time.Now().Unix()
	admin := datastruct.WebUser{
		Name:      "管理员",
		LoginName: "admin",
		Pwd:       "123@s678",
		RoleId:    adminLevelID,
		Token:     commondata.CommonDataInfo.UniqueId(),
		CreatedAt: now_time,
		UpdatedAt: now_time,
	}
	return []datastruct.WebUser{admin}
}*/

func initData(engine *xorm.Engine) {
	execStr := fmt.Sprintf("ALTER DATABASE %s CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;", conf.Server.DB_Name)
	_, err := engine.Exec(execStr)
	errhandle(err)

	// _, err = engine.Exec("ALTER TABLE user_info CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	// errhandle(err)

	// execStr = fmt.Sprintf("ALTER TABLE user_info AUTO_INCREMENT = %d", datastruct.UserIdStart)
	// _, err = engine.Exec(execStr)
	// errhandle(err)

	// _, err = engine.Exec("ALTER TABLE user_info CHANGE nick_name nick_name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	// errhandle(err)
}

func errhandle(err error) {
	if err != nil {
		log.Fatal("db error is %v", err.Error())
	}
}
