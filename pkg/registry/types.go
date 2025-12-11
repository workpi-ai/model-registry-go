package registry

import (
	"sync"

	"github.com/workpi-ai/go-utils/ghrelease"
)

type ProviderType string

type AuthType string

type APIFormat string

type Registry struct {
	Providers map[string]*Provider

	mu              sync.RWMutex
	configDir       string
	loader          *Loader
	updater         *ghrelease.Updater
	customProviders []*Provider
	stopChan        chan struct{}
	closeOnce       sync.Once
}

type Metadata struct {
	Version     string `json:"version"`
	LastCheckAt string `json:"last_check_at"`
}
