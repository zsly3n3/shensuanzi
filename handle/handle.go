package handle

import (
	"shensuanzi/cache"
	"shensuanzi/db"
	"sync"
)

type AppHandler struct {
	dbHandler    *db.DBHandler
	cacheHandler *cache.CACHEHandler
	payMutex     *sync.Mutex //读写互斥量
}

var appHandler *AppHandler

func CreateAppHandle() *AppHandler {
	appHandler = new(AppHandler)
	appHandler.cacheHandler = cache.CreateCACHEHandler()
	appHandler.dbHandler = db.CreateDBHandler(false)
	return appHandler
}

type WebHandler struct {
	dbHandler *db.DBHandler
}

var webHandler *WebHandler

func CreateWebHandle() *WebHandler {
	webHandler = new(WebHandler)
	webHandler.dbHandler = db.CreateDBHandler(true)
	return webHandler
}

func (handler *WebHandler) GetWebDBHandler() *db.DBHandler {
	return handler.dbHandler
}
