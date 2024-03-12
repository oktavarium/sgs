package sgs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/oktavarium/sgs/internal/broadcast"
)

const defaultUdpBufferSize = 64
const defaultSenderTimeout = 100 * time.Millisecond

type SGS struct {
	ctx           context.Context
	addr          string
	conn          *net.UDPConn
	senderTimeout time.Duration
	handlerFn     Handler
	senderFn      Sender
}

func NewServer(ctx context.Context,
	addr string,
	senderTimeout time.Duration,
	handler Handler,
	sender Sender,
) SGS {
	if senderTimeout == 0 {
		senderTimeout = defaultSenderTimeout
	}

	return SGS{
		ctx:           ctx,
		addr:          addr,
		senderTimeout: senderTimeout,
		handlerFn:     handler,
		senderFn:      sender,
	}
}

func (s *SGS) ListenAndServe() (err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return fmt.Errorf("parse addr: %w", err)
	}
	s.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("listen udp: %w", err)
	}

	wg := new(sync.WaitGroup)
	senderTicker := time.NewTicker(s.senderTimeout)
	broadcastSender := broadcast.NewSender(s.ctx, s.conn)
	wg.Add(1)
	go broadcastSender.Run()
	defer func() {
		defer senderTicker.Stop()
		wg.Wait()
		if err = s.conn.Close(); err != nil {
			err = fmt.Errorf("close udp socket: %w", err)
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			return fmt.Errorf("context is done: %w", s.ctx.Err())
		case <-senderTicker.C:
			clients, r := s.senderFn.Send()
			broadcastSender.Send(clients, r)
		default:
			var buf [defaultUdpBufferSize]byte
			if err := s.conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond)); err != nil {
				slog.Debug("set read deadline", "error", err)
				continue
			}
			size, client, err := s.conn.ReadFromUDPAddrPort(buf[:])
			if err != nil {
				slog.Debug("read from udp", "error", err)
				continue
			}
			go func() {
				wg.Add(1)
				defer wg.Done()
				r := bytes.NewBuffer(buf[:size])
				w := new(bytes.Buffer)
				s.handlerFn.ServeUDP(client, r, w)
				if w.Len() > 0 {
					data, err := io.ReadAll(w)
					if err != nil {
						slog.Debug("read from response", "error", err)
						return
					}
					_, err = s.conn.WriteToUDPAddrPort(data, client)
					if err != nil {
						slog.Debug("write to udp", "error", err)
						return
					}
				}
			}()
		}
	}
}
