package router

import (
	"fmt"
	"hello/auth"
	"hello/logging"
	"hello/modbus"
	"hello/model"
	"hello/opc"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// MsgFlags ...
var MsgFlags = map[int]string{
	1000: "ok",
	500:  "fail",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[500]
}

// Gin ...
type Gin struct {
	C *gin.Context
}

// Response ...
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  GetMsg(errCode),
		Data: data,
	})
	return
}

// var appG = app.Gin{C: c}

func getUser(c *gin.Context) {
	var appG = Gin{C: c}
	id := c.Param("id")
	name := c.Param("name")
	json := gin.H{
		"data": id,
		"name": name,
	}
	// c.JSON(http.StatusOK, json)
	appG.Response(http.StatusOK, 1000, json)
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

// 用户列表 共享变量（临界资源）
var userList []string

// 互斥锁
var mux sync.Mutex

//Index ...
func Index(c *gin.Context) {
	c.String(http.StatusOK, "当前参与抽奖的用户人数:%d", len(userList))
}

// ImportUsers ...
func ImportUsers(c *gin.Context) {
	strUsers := c.Query("users")
	users := strings.Split(strUsers, ",")
	// 在操作 全局变量 userList 之前加互斥锁，加完锁记得释放
	mux.Lock()
	defer mux.Unlock()
	// 统计当前已经在参加抽奖的用户数量
	currUserCount := len(userList)

	// 将页面提交的用户导入到 userList 中，参与抽奖
	for _, user := range users {
		user = strings.TrimSpace(user)
		if len(user) > 0 {
			userList = append(userList, user)
		}
	}
	// 统计当前总共参加抽奖人数
	userTotal := len(userList)
	c.String(http.StatusOK, "当前参与抽奖的用户数量:%d,导入的用户数量:%d", userTotal, (userTotal - currUserCount))
}

// GetLuckyUser ...
func GetLuckyUser(c *gin.Context) {
	var user string
	// 在操作 全局变量 userList 之前加互斥锁，加完锁记得释放
	mux.Lock()
	defer mux.Unlock()

	count := len(userList)
	if count > 1 {

		seed := time.Now().UnixNano()
		// 以随机数设置中奖用户, [0,count)中的随机值
		lotteryindex := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user = userList[lotteryindex]
		// 当前参与抽奖用户减 1
		userList = append(userList[0:lotteryindex], userList[lotteryindex+1:]...)
		c.String(http.StatusOK, "中奖用户为:%s，剩余用户数:%d", user, count-1)

	} else if count == 1 {
		user = userList[0]
		userList = userList[0:0] // 清空参与抽奖的用户列表
		c.String(http.StatusOK, "中奖用户为:%s，剩余用户数:%d", user, count-1)
	} else {
		c.String(http.StatusOK, "当前无参与抽奖的用户,请导入新的用户。")
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layout/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/views/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}

// InitRouter 初始化路由
func InitRouter() (r *gin.Engine) {
	r = gin.New()
	// r.HTMLRender = loadTemplates("./templates")
	// 自定义日志格式
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		return fmt.Sprintf("[%s] [%s] [%s] [%d] [%s]\n",
			param.TimeStamp.Format(`2006-01-02 15:04:05`), //请求时间
			//param.ClientIP,        //请求ip
			param.Method, //请求方法
			param.Path,   //请求路径
			//param.Request.Proto,   //协议
			param.StatusCode, //状态码
			//param.Latency,         //响应时间
			param.ErrorMessage, //错误信息
		)
	}))
	// 默认日志
	// r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
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
	userGroup := r.Group("/user")
	{
		// 首页
		userGroup.GET("/index", Index)
		// 导入用户
		userGroup.POST("/import", ImportUsers)
		// 抽奖
		userGroup.GET("/lucky", GetLuckyUser)
	}
	// 路由
	r.NoRoute(error404)
	// Any接受任何请求方法
	r.Any("/getorpost", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"page":    page,
			"name":    name,
			"message": message,
		})
	})
	r.GET("/users:name", getUser)
	r.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/index"
		r.HandleContext(c)

		logging.Debug("welcome to go lang")
	})
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
	// 获取url参数：welcome?firstname=hg&lastname=hyy
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
	// 初始化sqlite
	r.GET("/db", func(c *gin.Context) {

		model.UserDB()
		model.TestDB()
		c.JSON(http.StatusOK, gin.H{
			"code":    1000,
			"msg":     "db init  successful complete !",
			"success": true,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	var secrets = gin.H{
		"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
		"austin": gin.H{"email": "austin@example.com", "phone": "666"},
		"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
	}

	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets endpoint
	// hit "localhost:8080/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})
	return
}
