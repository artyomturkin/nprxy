package mw

import (
	"github.com/casbin/casbin"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// CasbinEnforcerConfig defines the config for CasbinEnforcer middleware.
	CasbinEnforcerConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Enforcer policy engine
		// Required.
		Enforcer *casbin.Enforcer

		// ParamGetters functions to source parameters for Casbin evaluation
		ParamGetters []func(echo.Context) interface{}
	}
)

// ValueFromContext get value from echo.Context
func ValueFromContext(key string) func(echo.Context) interface{} {
	return func(c echo.Context) interface{} {
		return c.Get(key)
	}
}

var (
	// defaultCasbinEnforcerConfig is the default CasbinEnforcer middleware config.
	defaultCasbinEnforcerConfig = CasbinEnforcerConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

// CasbinEnforcer returns a CasbinEnforcer middleware.
//
// For successful policy evaluation it calls the next handler.
// For failed evaluation, it sends "403 - Fobidden" response.
func CasbinEnforcer(ce *casbin.Enforcer, g ...func(echo.Context) interface{}) echo.MiddlewareFunc {
	c := defaultCasbinEnforcerConfig
	c.Enforcer = ce
	c.ParamGetters = g
	return CasbinEnforcerWithConfig(c)
}

// CasbinEnforcerWithConfig returns a CasbinEnforcer middleware with config.
// See `Middleware()`.
func CasbinEnforcerWithConfig(config CasbinEnforcerConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = defaultCasbinEnforcerConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) || check(c, config) {
				return next(c)
			}

			return echo.ErrForbidden
		}
	}
}

func check(c echo.Context, a CasbinEnforcerConfig) bool {
	p := make([]interface{}, 0)
	for _, f := range a.ParamGetters {
		p = append(p, f(c))
	}
	return a.Enforcer.Enforce(p...)
}
