package commondata

import (
	"shensuanzi/datastruct"
	"shensuanzi/db"
	"shensuanzi/log"
	"strconv"
	"sync"

	"github.com/holdno/snowFlakeByGo"
)

type CommonData struct {
	idWorker   *snowFlakeByGo.Worker
	serverInfo *ServerInfo
}

type ServerInfo struct {
	RWMutex    *sync.RWMutex //读写互斥量
	Version    string        //当前服务端版本号
	IsMaintain bool          //是否维护
}

var CommonDataInfo *CommonData
var once sync.Once

func Create(dbHandler *db.DBHandler) *CommonData {
	once.Do(func() {
		CommonDataInfo = new(CommonData)
		idWorker, err := snowFlakeByGo.NewWorker(0)
		if err != nil {
			log.Fatal("CreateCommonData err:%v", err.Error())
		}
		CommonDataInfo.idWorker = idWorker
		CommonDataInfo.serverInfo = createServerInfo(dbHandler)
	})
	return CommonDataInfo
}

func (data *CommonData) UniqueId() string {
	return strconv.FormatInt(data.idWorker.GetId(), 10)
}

func createServerInfo(dbHandler *db.DBHandler) *ServerInfo {
	server := new(ServerInfo)
	server.RWMutex = new(sync.RWMutex)
	db_data, code := dbHandler.GetServerInfo()
	if code != datastruct.NULLError {
		log.Fatal("GetServerInfo error from db")
		return nil
	}
	server.IsMaintain = db_data.IsMaintain
	server.Version = db_data.Version
	return server
}

func GetServerInfo() *ServerInfo {
	return CommonDataInfo.serverInfo
}
