package registry

import "fmt"

type Provider struct {
	Name        string            `yaml:"name" mapstructure:"name"`
	Type        ProviderType      `yaml:"type" mapstructure:"type"`
	AuthType    AuthType          `yaml:"auth_type" mapstructure:"auth_type"`
	APIKey      string            `yaml:"api_key" mapstructure:"api_key"`
	BaseURL     string            `yaml:"base_url" mapstructure:"base_url"`
	Description string            `yaml:"description" mapstructure:"description"`
	Models      map[string]*Model `yaml:"-" mapstructure:"models"`
}

func (p *Provider) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("provider name cannot be empty")
	}

	if p.AuthType == AuthTypeAPIKey && p.APIKey == "" {
		return fmt.Errorf("provider %s: api_key is required when auth_type is api_key", p.Name)
	}

	for _, model := range p.Models {
		if err := model.Validate(); err != nil {
			return fmt.Errorf("provider %s: %w", p.Name, err)
		}
	}

	return nil
}

func (p *Provider) Copy() *Provider {
	if p == nil {
		return nil
	}

	copied := &Provider{
		Name:        p.Name,
		Type:        p.Type,
		AuthType:    p.AuthType,
		APIKey:      p.APIKey,
		BaseURL:     p.BaseURL,
		Description: p.Description,
		Models:      make(map[string]*Model),
	}

	for name, model := range p.Models {
		copiedModel := model.Copy()
		copiedModel.Provider = copied
		copied.Models[name] = copiedModel
	}

	return copied
}

func (p *Provider) Merge(override *Provider) {
	SetIfNotZero(&p.Type, override.Type)
	SetIfNotZero(&p.AuthType, override.AuthType)
	SetIfNotZero(&p.APIKey, override.APIKey)
	SetIfNotZero(&p.BaseURL, override.BaseURL)
	SetIfNotZero(&p.Description, override.Description)

	if len(override.Models) > 0 {
		if p.Models == nil {
			p.Models = make(map[string]*Model)
		}
		for name, overrideModel := range override.Models {
			if existingModel, exists := p.Models[name]; exists {
				existingModel.Merge(overrideModel)
			} else {
				copiedModel := overrideModel.Copy()
				copiedModel.Provider = p
				p.Models[name] = copiedModel
			}
		}
	}
}
