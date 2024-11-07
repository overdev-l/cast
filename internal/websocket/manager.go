package websocket

import "sync"

type Manager struct {
	clients map[string]*Client
	mu      sync.RWMutex
}
