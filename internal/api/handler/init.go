package handler

import (
	"cast/config"
	"cast/internal/api/socket"
	"cast/internal/api/sse"

	"github.com/gin-gonic/gin"
)

type InitBody struct {
	LiveId int    `json:"live_id"`
	Plat   string `json:"plat"`
	Url    string `json:"url"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func InitHandler(c *gin.Context) {
	if config.LiveId != 0 {
		c.JSON(200, Response{
			Code: 0,
			Msg:  "live_id is default",
		})
		return
	}
	var body InitBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, Response{
			Code: 0,
			Msg:  "参数异常",
		})
		return
	}
	config.LiveId = body.LiveId
	config.Url = body.Url
	config.PlatFrom = body.Plat
	socket.ConnectWebSocket("ws://124.222.224.186:8800")
	socket.SendMessage("hello")
	c.JSON(200, Response{
		Code: 1,
		Msg:  "success",
	})
}

func StopHandler(c *gin.Context) {
	config.LiveId = 0
	config.Url = ""
	config.PlatFrom = ""
	sse.CloseSSE()
}

func SendMessage(c *gin.Context) {
	message := c.Query("message")
	sse.SendSSEMessage(message)
}
