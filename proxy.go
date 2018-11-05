package nprxy

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewProxy create proxy from configuration
func NewProxy(c ProxyConfig) (*http.Server, error) {
	switch c.Proto {
	case ProxyHTTP:
		return newHTTPProxy(c)
	default:
		return nil, fmt.Errorf("Unsupported protocol")
	}
}

func newHTTPProxy(c ProxyConfig) (*http.Server, error) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: c.UpstreamAddress}})))
	srv := &http.Server{
		Handler: e,
		Addr:    c.ListenAddress,
	}
	return srv, nil
}
