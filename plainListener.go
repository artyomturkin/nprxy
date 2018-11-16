package nprxy

import (
	"net"
)

func init() {
	listenerFactory["plain"] = buildPlainListener
}

func buildPlainListener(c ServiceConfig) (net.Listener, error) {
	return net.Listen("tcp", c.Listen.Address)
}
