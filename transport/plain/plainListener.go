package plain

import (
	"net"

	"github.com/artyomturkin/nprxy"
)

func init() {
	nprxy.ListenerFactory["plain"] = buildPlainListener
}

func buildPlainListener(c nprxy.ServiceConfig) (net.Listener, error) {
	return net.Listen("tcp", c.Listen.Address)
}
