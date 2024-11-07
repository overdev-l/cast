package message

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

var (
	manager *MessageManager
	once    sync.Once
)

// GetManager 获取单例消息管理器
func GetManager() *MessageManager {
	once.Do(func() {
		manager = &MessageManager{
			sseClients:  make(map[string]*SSEClient),
			messageChan: make(chan interface{}, 100),
			done:        make(chan struct{}),
		}
	})
	return manager
}

// ConnectWebSocket 连接WebSocket服务
func (m *MessageManager) ConnectWebSocket(url string) error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return err
	}

	m.wsConn = conn
	go m.handleWebSocketMessages()
	return nil
}

// handleWebSocketMessages 处理WebSocket消息
func (m *MessageManager) handleWebSocketMessages() {
	defer m.wsConn.Close()

	for {
		select {
		case <-m.done:
			return
		default:
			var message interface{}
			err := m.wsConn.ReadJSON(&message)
			if err != nil {
				log.Printf("WebSocket读取错误: %v", err)
				// 尝试重连
				m.reconnectWebSocket()
				continue
			}
			// 将消息转发到SSE
			m.messageChan <- message
		}
	}
}

// reconnectWebSocket WebSocket重连
func (m *MessageManager) reconnectWebSocket() {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * time.Duration(i+1))
		err := m.ConnectWebSocket(m.wsConn.URL().String())
		if err == nil {
			return
		}
	}
}

// BroadcastToSSE 广播消息到所有SSE客户端
func (m *MessageManager) BroadcastToSSE(event string, data interface{}) {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()

	for _, client := range m.sseClients {
		select {
		case client.MessageChan <- SSEMessage{
			Event: event,
			Data:  data,
		}:
		default:
			// 如果客户端消息队列满了，关闭连接
			client.Done <- struct{}{}
		}
	}
}

// Start 启动消息管理器
func (m *MessageManager) Start() {
	go func() {
		for {
			select {
			case msg := <-m.messageChan:
				m.BroadcastToSSE("message", msg)
			case <-m.done:
				return
			}
		}
	}()
}

// Stop 停止消息管理器
func (m *MessageManager) Stop() {
	close(m.done)
	if m.wsConn != nil {
		m.wsConn.Close()
	}
}
