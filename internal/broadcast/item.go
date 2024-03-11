package broadcast

import (
	"net/netip"
)

type sendItem struct {
	clients []netip.AddrPort
	data    []byte
}

func newSendItem(clients []netip.AddrPort, data []byte) sendItem {
	return sendItem{
		clients: clients,
		data:    data,
	}
}
