package main

import (
	// "app/event"

	// "app/tools"

	"net/http"
	"shensuanzi/commondata"
	"shensuanzi/conf"
	"shensuanzi/handle"
	"shensuanzi/routes/app"
	"shensuanzi/routes/web"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
)

//跨域
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Appversion, Apptoken, IsPc,Platform,Webtoken")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		//处理请求
		c.Next()
	}
}

func createData() {
	commondata.Create()

}

func registerRoutes(r *gin.Engine) {
	app_hanle := handle.CreateAppHandle()
	web_hanle := handle.CreateWebHandle()
	app.RegisterRoutes(r, app_hanle)
	web.RegisterRoutes(r, web_hanle)
}

func main() {
	r := gin.Default()
	var mode string
	switch conf.Common.Mode {
	case conf.Debug:
		mode = gin.DebugMode
	case conf.Release:
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	r.Use(cors())

	createData()
	registerRoutes(r)

	server := &http.Server{Addr: conf.Server.HttpServer, Handler: r}
	gracehttp.Serve(server)

	// r.Run(conf.Server.HttpServer) //listen and serve on 0.0.0.0:8080 asd
}
