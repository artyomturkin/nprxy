package nprxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// ProxyProto proxied protocol
type ProxyProto int

const (
	// ProxyHTTP proxies HTTP to specified endpoint
	ProxyHTTP ProxyProto = iota
	// ProxyTCP proxies TCP connection to specified address
	ProxyTCP
)

// ProxyDirection direction of proxied connection
type ProxyDirection int

const (
	// ProxyInbound forward traffic to impersonated service
	ProxyInbound ProxyDirection = iota
	// ProxyOutbound forward traffic to upstream dependecies of impersonated service
	ProxyOutbound
)

// ProxyConfig proxy service configuration
type ProxyConfig struct {
	Proto           ProxyProto
	Direction       ProxyDirection
	ListenAddress   string
	UpstreamAddress string
	TLS             *TLSConfig
}

// TLSConfig configuration for secure communication
type TLSConfig struct{}

// NewProxy create proxy from configuration
func NewProxy(ctx context.Context, c ProxyConfig) error {
	switch c.Proto {
	case ProxyHTTP:
		return newHTTPProxy(ctx, c)
	default:
		return fmt.Errorf("Unsupported protocol")
	}
}

func newHTTPProxy(ctx context.Context, c ProxyConfig) error {
	// Parse upstream address
	url, err := url.Parse(c.UpstreamAddress)
	if err != nil || (url.Scheme != "http" && url.Scheme != "https") {
		return fmt.Errorf("UpstreamAddress is not a valid url '%s': %v, or has unsupported scheme", c.UpstreamAddress, err)
	}

	// Setup listener for incoming connections
	listener, err := getListener(c.ListenAddress, c.TLS)
	if err != nil {
		return err
	}

	// Setup proxy HTTP handler and server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: url}})))

	srv := &http.Server{Handler: e}

	// Listen for termination in the background
	go func() {
		select {
		case <-ctx.Done():
			srv.Shutdown(ctx)
		}
	}()

	// Run server and handle requests in the background
	go srv.Serve(listener)

	// Continue
	return nil
}

func getListener(address string, c *TLSConfig) (net.Listener, error) {
	if c == nil {
		return net.Listen("tcp", address)
	}
	return nil, fmt.Errorf("Not implemented")
}
