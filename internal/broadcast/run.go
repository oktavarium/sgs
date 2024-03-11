package broadcast

import (
	"context"
	"log/slog"
	"net/netip"
	"sync"
	"time"
)

const deleteClientTimeout = 5 * time.Second

func (s *Sender) Run() {
	wg := sync.WaitGroup{}
	ticker := time.NewTicker(deleteClientTimeout)
	defer func() {
		defer ticker.Stop()
		for _, v := range s.clients {
			v.cancel()
			close(v.ch)
		}
		wg.Wait()
	}()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			for k, v := range s.clients {
				if time.Now().Sub(v.lastSend) > deleteClientTimeout {
					v.cancel()
					close(v.ch)
					delete(s.clients, k)
				}
			}
		case item := <-s.sendCh:
			for _, c := range item.clients {
				if cl, ok := s.clients[c]; ok {
					cl.lastSend = time.Now()
					cl.ch <- item.data
					continue
				}
				s.clients[c] = newClient(s.ctx)
				go func(ctx context.Context, client netip.AddrPort, ch chan []byte) {
					wg.Add(1)
					defer wg.Done()
					for {
						select {
						case <-ctx.Done():
							return
						case data := <-ch:
							_, err := s.conn.WriteToUDPAddrPort(data, client)
							if err != nil {
								slog.Debug("send data to client in broadcast mode", "error", err)
							}
						}
					}
				}(s.clients[c].ctx, c, s.clients[c].ch)
			}
		}
	}
}
