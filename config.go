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
	Name     string
	Listen   string
	Upstream string
	Grace    time.Duration

	// HTTP properties
	Timeout time.Duration
}
