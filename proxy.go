package nprxy

import (
	"context"
	"fmt"
	"net"
	"net/url"
)

// DialUpstream func to create conn to upstream service
type DialUpstream func(network, addr string) (net.Conn, error)

// Proxy forwards data from listener to upstream connection
type Proxy interface {
	Serve(ctx context.Context, Listener net.Listener, DialUpstream DialUpstream) error
}

// Proxy, Listener and Upstream factory
var (
	proxyFactory        = map[string]func(ServiceConfig) (Proxy, error){}
	listenerFactory     = map[string]func(ServiceConfig) (net.Listener, error){}
	upstreamDialFactory = map[string]func(ServiceConfig) (DialUpstream, error){}
)

// ProxyService create proxy and forward traffic
func ProxyService(ctx context.Context, c ServiceConfig) error {
	u, err := url.Parse(c.Upstream)
	if err != nil {
		return fmt.Errorf("failed to parse Upstream: %v", err)
	}

	// Create proxy with factory
	pf, ok := proxyFactory[u.Scheme]
	if !ok {
		return fmt.Errorf("unsupported upstream scheme %s", u.Scheme)
	}

	p, err := pf(c)
	if err != nil {
		return fmt.Errorf("failed to create proxy: %v", err)
	}

	// Create listener with factory
	if c.Listen.Kind == "" {
		c.Listen.Kind = "plain"
	}
	lf, ok := listenerFactory[c.Listen.Kind]
	if !ok {
		return fmt.Errorf("unsupported listener type %s", c.Listen.Kind)
	}

	l, err := lf(c)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}

	// Create upstream with factory
	udf, ok := upstreamDialFactory["plain"]
	if !ok {
		return fmt.Errorf("unsupported Upstream dialer type %s", "TODO")
	}

	ud, err := udf(c)
	if err != nil {
		return fmt.Errorf("failed to create upstream dialer: %v", err)
	}

	return p.Serve(ctx, l, ud)
}
