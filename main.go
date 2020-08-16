package main

import (
	"fmt"
	"hello/auth"
	"hello/handler"
	"hello/modbus"
	"hello/model"
	"hello/opc"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// import _ "github.com/jinzhu/gorm/dialects/mysql"
	// import _ "github.com/jinzhu/gorm/dialects/postgres"
	// import _ "github.com/jinzhu/gorm/dialects/sqlite"
	// import _ "github.com/jinzhu/gorm/dialects/mssql"
)

func getUser(c *gin.Context) {
	id := c.Param("id")
	name := c.Param("name")
	json := gin.H{
		"data": id,
		"name": name,
	}
	c.JSON(http.StatusOK, json)
}

func index(c *gin.Context) {
	c.Request.URL.Path = "/index"
	// r.HandleContext(c)

	log.Println("welcome to go lang")
}

func error404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

func test(response http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("base.html", "index.html")
	if err != nil {
		fmt.Println("parse index.html failed,err:", err)
		return
	}
	name := "tom"
	//tmpl.Execute(response,name)
	tmpl.ExecuteTemplate(response, "index.html", name)
}

func main() {

	// 模式
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	// 禁用控制台颜色
	gin.DisableConsoleColor()
	// 日志输出
	f, _ := os.Create("./log/gin.log")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(f)

	gin.DefaultWriter = io.MultiWriter(f)

	// 默认启动方式，包含 Logger、Recovery 中间件
	// r := gin.Default()
	// r.Use(gin.Recovery())
	// r.Use(gin.Logger())

	// 自定义启动
	r := gin.New()
	// r.Use(handler.Recover)

	http.HandleFunc("/", test)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	// 	log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	// }

	//自定义日志格式
	// r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

	// 	return fmt.Sprintf("[%s] [%s] [%s] [%d] [%s]\n",
	// 		param.TimeStamp.Format(`2006-01-02 15:04:05`), //请求时间
	// 		//param.ClientIP,        //请求ip
	// 		param.Method, //请求方法
	// 		param.Path,   //请求路径
	// 		//param.Request.Proto,   //协议
	// 		param.StatusCode, //状态码
	// 		//param.Latency,         //响应时间
	// 		param.ErrorMessage, //错误信息
	// 	)
	// }))

	// 静态文件
	r.Static("/assets", "./assets")
	// r.Static("/static", "./static")
	r.StaticFS("/static", http.Dir("./static"))
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	// HTML模板文件
	r.LoadHTMLGlob("templates/**/*")

	// 路由分组
	setting := r.Group("/setting")
	{
		setting.GET("/profile", auth.Profile)
	}

	admin := r.Group("/admin")
	{
		admin.POST("/signin", auth.Signin)
		admin.POST("/register", auth.Register)
		admin.GET("/signup", auth.Signup)
	}
	opcda := r.Group("/opc")
	{
		opcda.GET("/index", opc.Opcdaget)
		opcda.POST("/show", opc.Opcdapost)
	}
	mbs := r.Group("/modbus")
	{
		mbs.GET("/index", modbus.Modbusget)
		mbs.POST("/show", modbus.Modbuspost)
	}
	// 路由
	r.NoRoute(error404)
	r.Any("/test", index)
	r.GET("/user:name", getUser)
	r.GET("/", index)
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sign_in.html", gin.H{
			"title": "Main website",
		})

	})
	// cookie
	r.GET("/cookie", func(c *gin.Context) {

		cookie, err := c.Cookie("gin_cookie")

		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		fmt.Printf("Cookie value: %s \n", cookie)
	})
	// 获取get请求参数
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest") //如果没有则设置默认值
		lastname := c.Query("lastname")                   // 是 c.Request.URL.Query().Get("lastname") 的简写

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	// 获取第三方数据
	r.GET("/someDataFromReader", func(c *gin.Context) {
		response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})
	// src, des := `C:\Users\littl\Desktop\a.txt`, `D:\b.txt`
	// err := opc.CopyFile(src, des)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.GET("/db", func(c *gin.Context) {

		model.UserDB()
		model.TestDB()
		c.JSON(http.StatusOK, gin.H{
			"code":    1000,
			"msg":     "db init  successful complete !",
			"success": true,
		})
	})
	// fhh := handler.Persion{
	// 	Name: "fhh",
	// 	Age:  32,
	// }

	// in := 123
	// ft := 3.14
	// str := "golang"
	// arry := [...]int{1, 2, 3}
	// sli := []int{1, 2, 3}
	// mp := map[string]string{
	// 	"name": "fhh",
	// 	"age":  "32",
	// }
	// fmt.Println(in, str)

	// handler.Testreflect(&in)
	// handler.Testreflect(ft)
	// handler.Testreflect(&str)
	// handler.Testreflect(arry)
	// handler.Testreflect(sli)
	// handler.Testreflect(mp)
	// handler.Testreflect(&fhh)
	// fmt.Println(in, str)
	var wg sync.WaitGroup
	wg.Add(1)
	// go handler.Chantestprime()
	wg.Add(1)
	// go handler.Testlock()
	go handler.Readandwrite()
	s.ListenAndServe()
	wg.Wait()
}
