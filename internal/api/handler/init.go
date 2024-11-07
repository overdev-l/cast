package handler

import (
	"cast/config"
	"cast/pkg/message"
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
	manager := message.GetManager()
	if err := manager.ConnectWebSocket(body.Url); err != nil {
		c.JSON(200, Response{
			Code: 0,
			Msg:  "连接WebSocket失败: " + err.Error(),
		})
		return
	}
	manager.Start()
	c.JSON(200, gin.H{
		"code": 1,
		"data": gin.H{},
		"msg":  "success",
	})
}

func SSEHandler(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		c.JSON(400, Response{
			Code: 0,
			Msg:  "client_id is required",
		})
		return
	}

	// SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 创建SSE客户端
	client := &message.SSEClient{
		ID:          clientID,
		MessageChan: make(chan message.SSEMessage, 10),
		Done:        make(chan struct{}),
	}

	// 注册客户端
	manager := message.GetManager()
	manager.RegisterClient(client)
	defer manager.UnregisterClient(client)

	// 发送连接成功消息
	c.SSEvent("connect", gin.H{
		"client_id": clientID,
		"status":    "connected",
	})
	c.Writer.Flush()

	// 处理消息
	for {
		select {
		case msg := <-client.MessageChan:
			c.SSEvent(msg.Event, msg.Data)
			c.Writer.Flush()
		case <-client.Done:
			return
		case <-c.Request.Context().Done():
			return
		}
	}
}

func StopHandler(c *gin.Context) {
	manager := message.GetManager()
	manager.Stop()

	config.LiveId = 0
	config.Url = ""
	config.PlatFrom = ""

	c.JSON(200, Response{
		Code: 1,
		Msg:  "stopped",
	})
}
