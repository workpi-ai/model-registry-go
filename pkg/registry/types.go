package registry

type ProviderType string

type AuthType string

type APIFormat string

type Registry struct {
	Providers map[string]*Provider

	configDir       string
	loader          *Loader
	updater         *Updater
	customProviders []*Provider
	stopChan        chan struct{}
}

type Metadata struct {
	Version     string `json:"version"`
	LastCheckAt string `json:"last_check_at"`
}
