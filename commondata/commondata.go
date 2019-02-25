package commondata

import (
	"shensuanzi/log"
	"strconv"
	"sync"

	"github.com/holdno/snowFlakeByGo"
)

type CommonData struct {
	idWorker *snowFlakeByGo.Worker
}

var CommonDataInfo *CommonData
var once sync.Once

func Create() *CommonData {
	once.Do(func() {
		CommonDataInfo = new(CommonData)
		idWorker, err := snowFlakeByGo.NewWorker(0)
		if err != nil {
			log.Fatal("CreateCommonData err:%v", err.Error())
		}
		CommonDataInfo.idWorker = idWorker
	})
	return CommonDataInfo
}

func (data *CommonData) UniqueId() string {
	return strconv.FormatInt(data.idWorker.GetId(), 10)
}
