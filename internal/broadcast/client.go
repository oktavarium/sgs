package broadcast

import (
	"context"
	"time"
)

type client struct {
	ctx      context.Context
	cancel   context.CancelFunc
	ch       chan []byte
	lastSend time.Time
}

func newClient(ctx context.Context) *client {
	ctx, cancel := context.WithCancel(ctx)
	return &client{
		ctx:    ctx,
		cancel: cancel,
		ch:     make(chan []byte),
	}
}
