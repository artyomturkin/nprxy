package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	gohttp "net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/artyomturkin/nprxy"
	"github.com/artyomturkin/nprxy/middleware"
	"github.com/casbin/casbin"
	"github.com/labstack/echo/middleware"
	yaml "gopkg.in/yaml.v2"

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
	if !h.DisableLog {
		h.Middlewares = append(h.Middlewares, middleware.RequestID(), mw.Logrus())
	}
	if !h.DisableLog && c.HTTP.LogBody {
		h.Middlewares = append(h.Middlewares, middleware.BodyDump(mw.LogrusBodyLogger))
	}
	if h.Grace == 0 {
		h.Grace = 5 * time.Second // Set default grace period for shutdown
	}
	if h.Timeout == 0 {
		h.Timeout = 5 * time.Second // Set default timeout
	}
	if c.HTTP.Kind != "" {
		h.Middlewares = append(h.Middlewares, mw.OperationResolver(c.HTTP.Kind))
	}
	if c.HTTP.Authn != nil {
		if c.HTTP.Authn.Kind == "api-key" {
			yamlFile, err := ioutil.ReadFile(c.HTTP.Authn.Params["path"].(string))
			if err != nil {
				panic(fmt.Errorf("yamlFile.Get err   #%v ", err))
			}
			keys := map[string]string{}
			err = yaml.Unmarshal(yamlFile, keys)
			if err != nil {
				panic(fmt.Errorf("Unmarshal keys: %v", err))
			}

			h.Middlewares = append(h.Middlewares, mw.BCryptAPIKey(keys))
		}
	}
	if c.HTTP.Authz != nil {
		if c.HTTP.Authz.Kind == "casbin" {
			ce := casbin.NewEnforcer(c.HTTP.Authz.Params["model"].(string), c.HTTP.Authz.Params["policy"].(string))
			var p []func(echo.Context) interface{}
			for _, v := range c.HTTP.Authz.Params["parameters"].([]interface{}) {
				p = append(p, mw.ValueFromContext(v.(string)))
			}
			h.Middlewares = append(h.Middlewares, mw.CasbinEnforcer(ce, p...))
		}
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
