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
