package main

import (
	"fmt"
	"hello/docs"
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

	fmt.Println("+------------------------------------------------------------------+")
	fmt.Println("| Welcome to use Config tools                                      |")
	fmt.Println("| Code by hyy, latest update at 2020/08/19                         |")
	fmt.Println("| If you have any problem when using the tool                      |")
	fmt.Println("| Please submit issue at : https://github.com/hg-hyy/gg-gl |")
	fmt.Println("+------------------------------------------------------------------+")
}

// Help ...
func Help() {
	fmt.Println("+------------------------(-:---:)---------------------------------+")
	fmt.Println(`A: "-server" load "config.ini" and start as server`)
	fmt.Println(`B: "-client" load "config.ini" and start as client`)
	fmt.Println(`C: "for more details please read "README.md"`)
}

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @termsOfService http://swagger.io/terms/
func main() {

	// gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	// 	log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	// }
	docs.SwaggerInfo.Title = "Swagger API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "This is a config tools server."
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	defer func() {
		if err := recover(); err != nil {
			logging.Error(err)
		}
	}()

	// agrs := os.Args
	// if len(agrs) >= 2 && agrs[1] == "-h" {
	// 	Help()
	// } else {
	// 	panic("more args is needed")
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
