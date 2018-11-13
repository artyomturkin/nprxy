package nprxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/labstack/echo"
)

// Proxy forwards data from listener to upstream connection
type Proxy interface {
	Serve(ctx context.Context, Listener net.Listener, DialUpstream func(network, addr string) (net.Conn, error))
}

// HTTPProxy forwards HTTP requests to upstream service
type HTTPProxy struct {
	Upstream    *url.URL
	Grace       time.Duration
	Middlewares []echo.MiddlewareFunc
}

// Serve starts http server on listener, that uses connection from DialUpstream func to connect to upstream service and routes requests and response to and from upstream service
func (h *HTTPProxy) Serve(ctx context.Context, Listener net.Listener, DialUpstream func(network, addr string) (net.Conn, error)) error {
	if h.Grace == 0 {
		h.Grace = 5 * time.Second // Set default grace period for shutdown
	}

	r := httputil.NewSingleHostReverseProxy(h.Upstream)
	t := &http.Transport{
		Dial:    DialUpstream,
		DialTLS: DialUpstream,
	}
	r.Transport = t

	e := echo.New()
	e.Any("/*", echo.WrapHandler(r), h.Middlewares...)

	s := http.Server{
		Handler: e,
	}

	go func() {
		<-ctx.Done()

		c, cancel := context.WithTimeout(context.Background(), h.Grace)
		s.Shutdown(c)
		cancel()
	}()

	return s.Serve(Listener)
}
