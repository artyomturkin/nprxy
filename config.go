package nprxy

import (
	"net/url"
)

// Config for nprxy
type Config struct {
	Service struct {
		Name       string
		Endpoint   *url.URL
		Operations OperationNamerConfig
	}
	Listen string
}

// OperationNamerConfig configuration of OperationNamer
type OperationNamerConfig struct {
	Kind string
}

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
	UpstreamAddress *url.URL
}
