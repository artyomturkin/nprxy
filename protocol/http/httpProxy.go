package http

import (
	"context"
	"net"
	gohttp "net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/artyomturkin/nprxy"
	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

func init() {
	nprxy.ProxyFactory["http"] = buildHTTPProxy
	nprxy.ProxyFactory["https"] = buildHTTPProxy
}

func buildHTTPProxy(c nprxy.ServiceConfig) (nprxy.Proxy, error) {
	u, _ := url.Parse(c.Upstream)

	h := &httpProxy{
		Upstream:   u,
		Grace:      c.Grace,
		Timeout:    c.Timeout,
		DisableLog: c.DisableLog,
	}
	if h.Grace == 0 {
		h.Grace = 5 * time.Second // Set default grace period for shutdown
	}
	if h.Timeout == 0 {
		h.Timeout = 5 * time.Second // Set default timeout
	}
	return h, nil
}

// httpProxy forwards HTTP requests to upstream service
type httpProxy struct {
	Upstream    *url.URL
	Grace       time.Duration
	Timeout     time.Duration
	Middlewares []echo.MiddlewareFunc
	DisableLog  bool
}

// Serve starts http server on listener, that uses connection from DialUpstream func to connect to upstream service and routes requests and response to and from upstream service
func (h *httpProxy) Serve(ctx context.Context, Listener net.Listener, DialUpstream nprxy.DialUpstream) error {
	r := httputil.NewSingleHostReverseProxy(h.Upstream)
	t := &gohttp.Transport{
		Dial:    DialUpstream,
		DialTLS: DialUpstream,
		Proxy:   gohttp.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   h.Timeout,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	r.Transport = t

	rewriteHost := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Host = h.Upstream.Host
			return next(c)
		}
	}

	e := echo.New()
	mws := append(h.Middlewares, middleware.Secure(), rewriteHost)
	if !h.DisableLog {
		mws = append(mws, middleware.Logger())
	}
	e.Any("/*", echo.WrapHandler(r), mws...)

	s := gohttp.Server{
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
