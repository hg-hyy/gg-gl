package main

import (
	"fmt"
	"hello/opc"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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

//LoginForm struct
type LoginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func sign(c *gin.Context) {
	var form LoginForm
	// in this case proper binding will be automatically selected
	if c.ShouldBind(&form) == nil {

		if form.Email == "littleshenyun@outlook.com" && form.Password == "111111" {
			f, _ := c.FormFile("file")
			log.Println(f.Filename)
			c.SaveUploadedFile(f, "./upload")
			c.JSON(200, gin.H{
				"status":   "you are logged in",
				"Email":    form.Email,
				"password": form.Password,
				"success":  true})
		} else {
			c.JSON(401, gin.H{"status": "unauthorized"})
		}
	}
}
func login(c *gin.Context) {

	Email := c.PostForm("Email")
	password := c.PostForm("password")
	json := gin.H{
		"Email":    Email,
		"password": password,
	}
	c.JSON(http.StatusOK, json)
}
func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}
func main() {

	// 禁用控制台颜色
	gin.DisableConsoleColor()
	f, _ := os.Create("./log/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// 你的自定义格式
		return fmt.Sprintf(`[%s][%s][%s][%s][%s][%d][%s][%s]`,
			param.TimeStamp.Format(`2006-01-02 15:04:05`),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	v1 := r.Group("/v1")
	{
		v1.GET("/login", login)
	}

	// Simple group: v2
	v2 := r.Group("/v2")
	{
		v2.POST("/sign", sign)
	}

	// r.Static("/static", "./static")
	r.Static("/assets", "./assets")
	r.StaticFS("/static", http.Dir("./static"))
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	r.LoadHTMLGlob("templates/**/*")

	r.GET("/user:name", getUser)
	r.GET("/daget", opc.Opcdaget)
	r.POST("/dapost", opc.Opcdapost)

	r.POST("/sign", sign)

	r.GET("/", func(c *gin.Context) {

		c.Request.URL.Path = "/index"
		r.HandleContext(c)
		// c.JSON(http.StatusOK, gin.H{
		// 	"message": "hello,golang",
		// })
		log.Println("file.Filename")
	})
	r.GET("/login", func(c *gin.Context) {

		data := gin.H{
			"message": "login",
		}

		c.JSON(http.StatusOK, data)

	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Main website",
		})

	})

	r.GET("/cookie", func(c *gin.Context) {

		cookie, err := c.Cookie("gin_cookie")

		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		fmt.Printf("Cookie value: %s \n", cookie)
	})
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest") //如果没有则设置默认值
		lastname := c.Query("lastname")                   // 是 c.Request.URL.Query().Get("lastname") 的简写

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	// src, des := `C:\Users\littl\Desktop\a.txt`, `D:\b.txt`

	// err := opc.CopyFile(src, des)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	s.ListenAndServe()
}
