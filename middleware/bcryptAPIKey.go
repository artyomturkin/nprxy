package mw

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/bcrypt"
)

type (
	// BCryptAPIKeyConfig defines the config for BCryptAPIKey middleware.
	BCryptAPIKeyConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Keys system - hashed key pairs to check against
		Keys map[string]string

		// ClientHeader header that carries system name
		ClientHeader string

		// KeyHeader header that carries system key
		KeyHeader string

		// ContextKey key to output client if authenticated
		ContextKey string
	}
)

var (
	// defaultBCryptAPIKey is the default BCryptAPIKey middleware config.
	defaultBCryptAPIKey = BCryptAPIKeyConfig{
		Skipper:      middleware.DefaultSkipper,
		Keys:         map[string]string{},
		ClientHeader: "X-NPRXY-Client",
		KeyHeader:    "X-NPRXY-Key",
		ContextKey:   "client",
	}
)

// BCryptAPIKey returns a BCryptAPIKey middleware.
func BCryptAPIKey(keys map[string]string) echo.MiddlewareFunc {
	c := defaultBCryptAPIKey
	c.Keys = keys
	return BCryptAPIKeyWithConfig(c)
}

// BCryptAPIKeyWithConfig returns a BCryptAPIKey middleware with config.
// See `Middleware()`.
func BCryptAPIKeyWithConfig(config BCryptAPIKeyConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = defaultBCryptAPIKey.Skipper
	}

	if config.Keys == nil {
		config.Keys = map[string]string{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			system := c.Request().Header.Get(config.ClientHeader)
			key := c.Request().Header.Get(config.KeyHeader)
			hash, ok := config.Keys[system]

			if ok {
				err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(key))
				if err == nil {
					c.Set(config.ContextKey, system)
					return next(c)
				}
			}

			return echo.ErrUnauthorized
		}
	}
}
