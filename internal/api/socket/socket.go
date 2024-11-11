package socket

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	conn      *websocket.Conn
	connMutex sync.Mutex
	isActive  bool
)

// ConnectWebSocket 连接到 WebSocket 服务器
func ConnectWebSocket(url string) error {
	connMutex.Lock()
	defer connMutex.Unlock()

	if isActive {
		return fmt.Errorf("websocket connection already exists")
	}

	// 创建 websocket 连接
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("dial error: %v", err)
	}

	conn = c
	isActive = true

	// 启动消息处理协程
	go handleMessages()

	return nil
}

// handleMessages 处理接收到的消息
func handleMessages() {
	defer func() {
		connMutex.Lock()
		isActive = false
		conn.Close()
		connMutex.Unlock()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			return
		}

		// 处理接收到的消息
		handleReceivedMessage(message)
	}
}

// handleReceivedMessage 处理收到的消息
func handleReceivedMessage(message []byte) {
	// 这里处理接收到的消息
	log.Printf("Received message: %s", string(message))
}

// CloseConnection 关闭 WebSocket 连接
func CloseConnection() error {
	connMutex.Lock()
	defer connMutex.Unlock()

	if !isActive {
		return fmt.Errorf("no active connection")
	}

	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return fmt.Errorf("write close error: %v", err)
	}

	return conn.Close()
}

// IsConnected 检查连接状态
func IsConnected() bool {
	connMutex.Lock()
	defer connMutex.Unlock()
	return isActive
}
