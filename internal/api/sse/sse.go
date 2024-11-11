package sse

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

var (
	messageChan chan string // 用于发送消息的通道
	closeChan   chan bool   // 用于关闭连接的通道
	isConnected bool        // 连接状态
)

func init() {
	messageChan = make(chan string)
	closeChan = make(chan bool)
	isConnected = false
}

func SSEHandler(c *gin.Context) {
	// 如果已经有连接，拒绝新的连接
	if isConnected {
		c.JSON(400, gin.H{"error": "connection already exists"})
		return
	}

	// 设置 SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	isConnected = true
	defer func() { isConnected = false }()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-messageChan:
			c.SSEvent("message", msg)
			return true
		case <-closeChan:
			return false
		}
	})
}

// SendSSEMessage 发送消息到 SSE 客户端
func SendSSEMessage(message string) error {
	if !isConnected {
		return fmt.Errorf("no active SSE connection")
	}
	messageChan <- message
	return nil
}

// CloseSSE 关闭 SSE 连接
func CloseSSE() {
	if isConnected {
		closeChan <- true
	}
}
