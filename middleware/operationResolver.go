package mw

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// OperationResolverConfig defines the config for OperationResolver middleware.
	OperationResolverConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// ContextKey key to output client if authenticated
		ContextKey string

		// Kind switches operation resolution logic. Supported: "soap"
		Kind string
	}
)

var (
	// defaultOperationResolver is the default OperationResolver middleware config.
	defaultOperationResolver = OperationResolverConfig{
		Skipper:    middleware.DefaultSkipper,
		ContextKey: "operation",
	}
)

// OperationResolver returns a OperationResolver middleware.
func OperationResolver(kind string) echo.MiddlewareFunc {
	c := defaultOperationResolver
	c.Kind = kind
	return OperationResolverWithConfig(c)
}

// OperationResolverWithConfig returns a OperationResolver middleware with config.
// See `Middleware()`.
func OperationResolverWithConfig(config OperationResolverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = defaultOperationResolver.Skipper
	}

	resolvers := map[string]func(r *http.Request) (string, error){}

	resolvers["soap"] = func(r *http.Request) (string, error) {
		if op := r.Header.Get("SOAPAction"); op != "" {
			return op, nil
		}

		return "", fmt.Errorf("SOAPAction not set")
	}

	if resolver, ok := resolvers[config.Kind]; ok {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if config.Skipper(c) {
					return next(c)
				}

				op, err := resolver(c.Request())
				if err != nil {
					return echo.ErrBadRequest
				}

				c.Set(config.ContextKey, op)
				return next(c)
			}
		}
	}
	panic(fmt.Errorf("unsupported operation resolver kind"))
}
