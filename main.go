package main

import (
	"fmt"
	"hello/logging"
	"hello/router"
	"hello/setting"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {

	setting.Setup()
	logging.Setup()
}

func main() {

	// gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	// 	log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	// }

	gin.SetMode(setting.ServerSetting.RunMode)
	routeHandler := router.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%s", setting.ServerSetting.HTTPPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routeHandler,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logging.Info("start http server listening", endPoint)

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go handler.Testlock()
	// go handler.Chantestprime()
	// go handler.Readandwrite()
	// wg.Wait()
	server.ListenAndServe()
}
