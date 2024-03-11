package sgs

import (
	"io"
	"net/netip"
)

type Sender interface {
	Send() ([]netip.AddrPort, io.Reader)
}

type SenderFunc func() ([]netip.AddrPort, io.Reader)

func (f SenderFunc) Send() ([]netip.AddrPort, io.Reader) {
	return f()
}
