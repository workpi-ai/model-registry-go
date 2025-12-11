package registry

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	embed "github.com/workpi-ai/model-registry-go"
	"gopkg.in/yaml.v3"
)

const (
	minModelPathSegments = 3
)

type Loader struct {
	configDir string
}

func NewLoader(configDir string) *Loader {
	return &Loader{configDir: configDir}
}

func (l *Loader) Load() (map[string]*Provider, error) {
	var fsys fs.FS

	localPath := filepath.Join(l.configDir, providersDir)
	if stat, err := os.Stat(localPath); err == nil && stat.IsDir() {
		fsys = os.DirFS(localPath)
	} else {
		embedFS, err := embed.GetFS()
		if err != nil {
			return nil, fmt.Errorf("failed to load embedded data: %w", err)
		}
		fsys = embedFS
	}

	providers := make(map[string]*Provider)
	return providers, l.parseFS(providers, fsys)
}

func (l *Loader) parseFS(providers map[string]*Provider, fsys fs.FS) error {
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, yamlExt) {
			return nil
		}

		if !strings.HasSuffix(path, providerYAML) {
			return nil
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		return l.parseProvider(providers, path, data)
	})
	if err != nil {
		return err
	}

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, yamlExt) {
			return nil
		}

		if strings.HasSuffix(path, providerYAML) {
			return nil
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		return l.parseModel(providers, path, data)
	})

	return err
}

func (l *Loader) parseProvider(providers map[string]*Provider, path string, data []byte) error {
	var provider Provider
	if err := yaml.Unmarshal(data, &provider); err != nil {
		return fmt.Errorf("parse provider %s: %w", path, err)
	}

	provider.Models = make(map[string]*Model)
	providers[provider.Name] = &provider
	return nil
}

func (l *Loader) parseModel(providers map[string]*Provider, path string, data []byte) error {
	var model Model
	if err := yaml.Unmarshal(data, &model); err != nil {
		return fmt.Errorf("parse model %s: %w", path, err)
	}

	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) < minModelPathSegments {
		return fmt.Errorf("invalid model path: %s", path)
	}

	providerName := parts[0]

	provider, ok := providers[providerName]
	if !ok {
		return nil
	}

	model.Provider = provider
	provider.Models[model.Name] = &model

	return nil
}


