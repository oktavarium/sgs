package broadcast

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/netip"
)

type Sender struct {
	ctx     context.Context
	conn    *net.UDPConn
	sendCh  chan sendItem
	clients map[netip.AddrPort]*client
}

func NewSender(ctx context.Context, conn *net.UDPConn) Sender {
	return Sender{
		ctx:     ctx,
		conn:    conn,
		sendCh:  make(chan sendItem),
		clients: make(map[netip.AddrPort]*client),
	}
}

func (s *Sender) Send(r io.Reader, addrs ...netip.AddrPort) {
	if addrs == nil {
		slog.Debug("no clients for broadcast")
		return
	}
	data, err := io.ReadAll(r)
	if err != nil || data == nil {
		slog.Debug("reading from send reader", "error", err)
		return
	}
	s.sendCh <- newSendItem(addrs, data)
}
