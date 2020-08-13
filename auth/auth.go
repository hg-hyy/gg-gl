package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//LoginForm struct
type LoginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

//Sign    login
func Sign(c *gin.Context) {
	var form LoginForm
	// in this case proper binding will be automatically selected
	if c.ShouldBind(&form) == nil {

		if form.Email == "littleshenyun@outlook.com" && form.Password == "111111" {

			c.JSON(http.StatusOK, gin.H{
				"status":   "you are logged in",
				"Email":    form.Email,
				"password": form.Password,
				"success":  true})
			log.Println("登录成功")
		} else {
			c.JSON(401, gin.H{"status": "unauthorized"})
		}
	}
}

// Register  Register
func Register(c *gin.Context) {
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

// Profile account
func Profile(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", nil)
}
