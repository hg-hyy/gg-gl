package modbus

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Modbus ...
type Modbus struct {
	DevID    int    `json:"dev_id"`
	CollType string `json:"Coll_Type"`
	TCP      TCP    `json:"TCP"`
	RTU      RTU    `json:"RTU"`
}

// TCP setting
type TCP struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// RTU setting
type RTU struct {
	Serial  string `json:"serial"`
	Baud    int    `json:"baud"`
	DataBit int    `json:"data_bit"`
	StopBit int    `json:"stop_bit"`
	Parity  string `json:"parity"`
}

// Modbusget ...
func Modbusget(c *gin.Context) {
	c.HTML(http.StatusOK, "modbus_index.html", nil)

}

// Modbuspost ...
func Modbuspost(c *gin.Context) {
	c.HTML(http.StatusOK, "modbus_show.html", gin.H{
		"code":    1000,
		"msg":     "config success",
		"success": true,
	})

}
