package websocket

import (
	"cast/internal/websocket"
	"cast/pkg/models"
	"context"
	"sync"

	"golang.org/x/net/websocket"
)

type Client struct {
	url        string
	name       string
	conn       *websocket.Conn
	message    chan models.Message
	ctx        context.Context
	cancelFunc context.CancelCauseFunc
	mu         sync.Mutex
}
