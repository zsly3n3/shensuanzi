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
	// appHandler.cacheHandler = cache.CreateCACHEHandler()
	// appHandler.dbHandler = db.CreateDBHandler()
	return appHandler
}

type WebHandler struct {
	dbHandler    *db.DBHandler
	cacheHandler *cache.CACHEHandler
}

var webHandler *WebHandler

func CreateWebHandle() *WebHandler {
	webHandler = new(WebHandler)
	// webHandler.cacheHandler = cache.CreateCACHEHandler()
	// webHandler.dbHandler = db.CreateDBHandler()
	return webHandler
}
