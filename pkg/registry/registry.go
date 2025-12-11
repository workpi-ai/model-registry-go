package registry

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"time"
)

const (
	DefaultCheckInterval = 1 * time.Hour
	defaultDirPerm       = 0755
	defaultFilePerm      = 0644
)

type Options struct {
	ConfigDir     string
	AutoUpdate    bool
	CheckInterval time.Duration
	Providers     []*Provider
}

func New(opts Options) (*Registry, error) {
	if opts.ConfigDir == "" {
		return nil, fmt.Errorf("ConfigDir is required")
	}
	if opts.CheckInterval == 0 {
		opts.CheckInterval = DefaultCheckInterval
	}

	if err := os.MkdirAll(opts.ConfigDir, defaultDirPerm); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	reg := &Registry{
		Providers:       make(map[string]*Provider),
		configDir:       opts.ConfigDir,
		loader:          NewLoader(opts.ConfigDir),
		customProviders: opts.Providers,
		stopChan:        make(chan struct{}),
	}

	updater, err := NewUpdater(opts.ConfigDir)
	if err != nil {
		return nil, fmt.Errorf("create updater: %w", err)
	}
	reg.updater = updater
	
	if err := reg.reload(); err != nil {
		return nil, err
	}

	if opts.AutoUpdate {
		go reg.autoUpdateLoop(opts.CheckInterval)
	}

	return reg, nil
}

func (r *Registry) Provider(name string) *Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Providers[name]
}

func (r *Registry) ProviderList() []*Provider {
	r.mu.RLock()
	providers := make([]*Provider, 0, len(r.Providers))
	for _, p := range r.Providers {
		providers = append(providers, p)
	}
	r.mu.RUnlock()

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Name < providers[j].Name
	})
	return providers
}

func (r *Registry) Model(provider, modelName string) *Model {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.Providers[provider]
	if !ok {
		return nil
	}

	return p.Models[modelName]
}

func (r *Registry) reload() error {
	newProviders, err := r.loader.Load()
	if err != nil {
		return err
	}

	if err := r.mergeCustomProviders(newProviders); err != nil {
		return err
	}

	r.mu.Lock()
	r.Providers = newProviders
	r.mu.Unlock()

	return nil
}

func (r *Registry) mergeCustomProviders(providers map[string]*Provider) error {
	for _, customProvider := range r.customProviders {
		existingProvider := providers[customProvider.Name]
		if existingProvider == nil {
			providers[customProvider.Name] = customProvider.Copy()
		} else {
			existingProvider.Merge(customProvider)
		}

		if err := providers[customProvider.Name].Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) autoUpdateLoop(interval time.Duration) {
	if err := r.updater.Update(); err == nil {
		_ = r.reload()
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.updater.Update(); err != nil {
				slog.Error("failed to update registry", "error", err)
				continue
			}

			if err := r.reload(); err != nil {
				slog.Error("failed to reload registry", "error", err)
				continue
			}
		case <-r.stopChan:
			return
		}
	}
}

func (r *Registry) Close() error {
	r.closeOnce.Do(func() {
		close(r.stopChan)
	})
	return nil
}

func (r *Registry) ForceUpdate() error {
	if err := r.updater.Update(); err != nil {
		return err
	}

	return r.reload()
}
