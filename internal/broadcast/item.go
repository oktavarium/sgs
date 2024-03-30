package broadcast

import (
	"net/netip"
)

type sendItem struct {
	addrs []netip.AddrPort
	data  []byte
}

func newSendItem(addrs []netip.AddrPort, data []byte) sendItem {
	return sendItem{
		addrs: addrs,
		data:  data,
	}
}
