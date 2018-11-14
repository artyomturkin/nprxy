package nprxy

import (
	"context"
	"net"
)

// DialUpstream func to create conn to upstream service
type DialUpstream func(network, addr string) (net.Conn, error)

// Proxy forwards data from listener to upstream connection
type Proxy interface {
	Serve(ctx context.Context, Listener net.Listener, DialUpstream DialUpstream) error
}

// Proxy factory
type buildProxy func(ServiceConfig) (Proxy, error)

var proxyBuilders = map[string]buildProxy{}
