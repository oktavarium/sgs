package sgs

import (
	"io"
	"net/netip"
)

type Handler interface {
	ServeUDP(netip.AddrPort, io.Reader, io.Writer)
}

type HandlerFunc func(netip.AddrPort, io.Reader, io.Writer)

func (f HandlerFunc) ServeUDP(addr netip.AddrPort, r io.Reader, w io.Writer) {
	f(addr, r, w)
}

func (s *SGS) serve(addr netip.AddrPort, r io.Reader, w io.Writer) {

}
