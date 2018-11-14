package nprxy

import (
	"context"
	"net"
)

// Proxy forwards data from listener to upstream connection
type Proxy interface {
	Serve(ctx context.Context, Listener net.Listener, DialUpstream func(network, addr string) (net.Conn, error))
}

var proxyBuilders = map[string]proxyBuilder{}

type proxyBuilder interface {
	Build(ServiceConfig) (*Proxy, error)
}
