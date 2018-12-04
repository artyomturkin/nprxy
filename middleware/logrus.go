package mw

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// LogrusConfig defines the config for Logrus middleware.
	LogrusConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
	}
)

var (
	// defaultLogrusConfig is the default LogrusConfig middleware config.
	defaultLogrusConfig = LogrusConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

func Logrus() echo.MiddlewareFunc {
	return LogrusWithConfig(defaultLogrusConfig)
}

func LogrusWithConfig(config LogrusConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = defaultLogrusConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			p := req.URL.Path
			if p == "" {
				p = "/"
			}

			bytesIn := req.Header.Get(echo.HeaderContentLength)
			if bytesIn == "" {
				bytesIn = "0"
			}

			logrus.WithFields(map[string]interface{}{
				"time_rfc3339":  time.Now().Format(time.RFC3339),
				"request_id":    res.Header().Get(echo.HeaderXRequestID),
				"remote_ip":     c.RealIP(),
				"host":          req.Host,
				"uri":           req.RequestURI,
				"method":        req.Method,
				"path":          p,
				"referer":       req.Referer(),
				"user_agent":    req.UserAgent(),
				"status":        res.Status,
				"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
				"latency_human": stop.Sub(start).String(),
				"bytes_in":      bytesIn,
				"bytes_out":     strconv.FormatInt(res.Size, 10),
			}).Info("Handled request")

			return nil
		}
	}
}

func LogrusBodyLogger(c echo.Context, reqB []byte, resB []byte) {
	res := c.Response()

	if len(reqB) > 0 {
		logrus.WithFields(map[string]interface{}{
			"request_id": res.Header().Get(echo.HeaderXRequestID),
			"body":       fmt.Sprintf("%.20000s", string(reqB)),
		}).Info("Request body")
	}

	if len(resB) > 0 {
		logrus.WithFields(map[string]interface{}{
			"request_id": res.Header().Get(echo.HeaderXRequestID),
			"body":       fmt.Sprintf("%.20000s", string(resB)),
		}).Info("Response body")
	}
}
