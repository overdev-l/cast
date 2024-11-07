// pkg/message/client.go
package message

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

// SSEMessage SSE消息结构
type SSEMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	ID    string      `json:"id,omitempty"`
}

// SSEClient SSE客户端结构
type SSEClient struct {
	ID          string          // 客户端唯一标识
	MessageChan chan SSEMessage // 消息通道
	Done        chan struct{}   // 关闭信号
	LastEventID string          // 最后一个事件ID
	CreatedAt   time.Time       // 创建时间
}

// MessageManager 消息管理器
type MessageManager struct {
	wsConn       *websocket.Conn
	sseClients   map[string]*SSEClient
	clientsMutex sync.RWMutex
	messageChan  chan interface{}
	done         chan struct{}
}

// RegisterClient 注册SSE客户端
func (m *MessageManager) RegisterClient(client *SSEClient) {
	m.clientsMutex.Lock()
	defer m.clientsMutex.Unlock()

	client.CreatedAt = time.Now()
	m.sseClients[client.ID] = client
}

// UnregisterClient 注销SSE客户端
func (m *MessageManager) UnregisterClient(client *SSEClient) {
	m.clientsMutex.Lock()
	defer m.clientsMutex.Unlock()

	if _, exists := m.sseClients[client.ID]; exists {
		close(client.Done)
		delete(m.sseClients, client.ID)
	}
}

// GetClient 获取SSE客户端
func (m *MessageManager) GetClient(clientID string) (*SSEClient, bool) {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()

	client, exists := m.sseClients[clientID]
	return client, exists
}

// SendToClient 发送消息到指定客户端
func (m *MessageManager) SendToClient(clientID string, event string, data interface{}) bool {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()

	if client, exists := m.sseClients[clientID]; exists {
		select {
		case client.MessageChan <- SSEMessage{
			Event: event,
			Data:  data,
			ID:    time.Now().String(),
		}:
			return true
		default:
			return false
		}
	}
	return false
}

// CleanInactiveClients 清理不活跃的客户端
func (m *MessageManager) CleanInactiveClients(timeout time.Duration) {
	m.clientsMutex.Lock()
	defer m.clientsMutex.Unlock()

	now := time.Now()
	for id, client := range m.sseClients {
		if now.Sub(client.CreatedAt) > timeout {
			close(client.Done)
			delete(m.sseClients, id)
		}
	}
}

// StartCleanup 启动清理任务
func (m *MessageManager) StartCleanup() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.CleanInactiveClients(30 * time.Minute)
			case <-m.done:
				return
			}
		}
	}()
}

// GetConnectedClients 获取已连接的客户端数量
func (m *MessageManager) GetConnectedClients() int {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()
	return len(m.sseClients)
}

// GetClientStatus 获取客户端状态
func (m *MessageManager) GetClientStatus(clientID string) map[string]interface{} {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()

	if client, exists := m.sseClients[clientID]; exists {
		return map[string]interface{}{
			"id":           client.ID,
			"connected_at": client.CreatedAt,
			"last_event":   client.LastEventID,
			"active":       true,
		}
	}

	return map[string]interface{}{
		"id":     clientID,
		"active": false,
	}
}
