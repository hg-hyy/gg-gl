package auth

import (
	"fmt"
	"hello/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//LoginForm struct
type LoginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// RegisterForm sign up
type RegisterForm struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

//Sign    login
func Signin(c *gin.Context) {
	var form LoginForm
	// in this case proper binding will be automatically selected
	if c.ShouldBind(&form) == nil {

		db, err := gorm.Open("sqlite3", "./model/hello.sqlite")
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()
		ok := db.Where(&model.User{Email: form.Email, Password: form.Password}).First(&model.User{})
		if ok.RecordNotFound() {

			c.JSON(401, gin.H{
				"code":    500,
				"msg":     "username or password is not right,please try again",
				"success": false,
			})
			log.Println("登录失败")
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    1000,
				"msg":     "you are success login !",
				"success": true,
			})
		}
	}
}

// Register  Register
func Register(c *gin.Context) {

	var form RegisterForm
	// in this case proper binding will be automatically selected
	if c.ShouldBind(&form) == nil {

		if form.Email != "" && form.Password != "" {
			user := model.User{
				Username: form.Username,
				Email:    form.Email,
				Password: form.Password,
			}
			db, err := gorm.Open("sqlite3", "./model/hello.sqlite")
			if err != nil {
				fmt.Println(err)
			}
			defer db.Close()
			ok := db.Where("Username = ?", form.Username).First(&model.User{})
			if ok.RecordNotFound() {
				db.Create(&user)
				c.JSON(200, gin.H{
					"code":    1000,
					"msg":     "register success !",
					"success": true,
				})
			} else {
				c.JSON(200, gin.H{
					"code":    500,
					"msg":     "register fail ! username already exists",
					"success": false,
				})
			}

		} else {
			c.JSON(401, gin.H{"status": "unauthorized"})
		}
	}
}

// Signup 注册
func Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "sign_up.html", nil)

}

// Profile account
func Profile(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", nil)
}
