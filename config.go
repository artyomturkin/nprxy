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

	// HTTP properties
	Timeout time.Duration
}

// ListenerConfig configuration of inbound channel
type ListenerConfig struct {
	Address string
	Kind    string
	TLSCert string `json:"tls_cert"`
	TLSKey  string `json:"tls_key"`
}
