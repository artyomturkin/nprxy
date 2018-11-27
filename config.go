package nprxy

import (
	"time"
)

// Config for nprxy
type Config struct {
	Services []ServiceConfig
}

// ServiceConfig general service configuration
type ServiceConfig struct {
	Name       string
	Listen     ListenerConfig
	Upstream   string
	Grace      time.Duration
	DisableLog bool

	// RPC properties
	Timeout time.Duration

	HTTP HTTPConfig
}

// ListenerConfig configuration of inbound channel
type ListenerConfig struct {
	Address string
	Kind    string
	TLSCert string `json:"tls_cert"`
	TLSKey  string `json:"tls_key"`
}

// HTTPConfig configuration for HTTP protocol
type HTTPConfig struct {
	Kind  string
	Authn *Parameters
	Authz *Parameters
}

// Parameters of config
type Parameters struct {
	Kind   string
	Params map[string]interface{}
}
