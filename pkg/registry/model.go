package registry

import "fmt"

type Model struct {
	Name     string    `yaml:"name" mapstructure:"name"`
	Provider *Provider `yaml:"-" mapstructure:"-"`
	APIs     APIs      `yaml:"apis" mapstructure:"apis"`
}

func (m *Model) Copy() *Model {
	if m == nil {
		return nil
	}

	model := &Model{
		Name:     m.Name,
		Provider: m.Provider,
		APIs: APIs{
			ChatCompletion: m.APIs.ChatCompletion.Copy(),
		},
	}

	return model
}

func (m *Model) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	if m.APIs.ChatCompletion != nil {
		if m.APIs.ChatCompletion.Context.MaxInput <= 0 {
			return fmt.Errorf("model %s: max_input must be positive", m.Name)
		}
		if m.APIs.ChatCompletion.Context.MaxOutput <= 0 {
			return fmt.Errorf("model %s: max_output must be positive", m.Name)
		}
		if m.APIs.ChatCompletion.Context.MaxOutput > m.APIs.ChatCompletion.Context.MaxInput {
			return fmt.Errorf("model %s: max_output (%d) cannot exceed max_input (%d)", 
				m.Name, 
				m.APIs.ChatCompletion.Context.MaxOutput, 
				m.APIs.ChatCompletion.Context.MaxInput)
		}
		if m.APIs.ChatCompletion.Parameters.MaxTokens <= 0 {
			return fmt.Errorf("model %s: max_tokens must be positive", m.Name)
		}
	}

	return nil
}

func (m *Model) Merge(override *Model) {
	SetIfNotZero(&m.Name, override.Name)

	if override.APIs.ChatCompletion != nil {
		if m.APIs.ChatCompletion == nil {
			m.APIs.ChatCompletion = override.APIs.ChatCompletion.Copy()
		} else {
			m.APIs.ChatCompletion.Merge(override.APIs.ChatCompletion)
		}
	}
}

type APIs struct {
	ChatCompletion *ChatCompletion `yaml:"chat_completion" mapstructure:"chat_completion"`
}

type ChatCompletion struct {
	APIFormat  APIFormat  `yaml:"api_format" mapstructure:"api_format"`
	Endpoint   string     `yaml:"endpoint" mapstructure:"endpoint"`
	Context    Context    `yaml:"context" mapstructure:"context"`
	Features   Features   `yaml:"features" mapstructure:"features"`
	Parameters Parameters `yaml:"parameters" mapstructure:"parameters"`
}

type Parameters struct {
	Temperature float64        `yaml:"temperature" mapstructure:"temperature"`
	MaxTokens   int            `yaml:"max_tokens" mapstructure:"max_tokens"`
	Extra       map[string]any `yaml:",inline" mapstructure:",remain"`
}

func (p Parameters) Copy() Parameters {
	copied := Parameters{
		Temperature: p.Temperature,
		MaxTokens:   p.MaxTokens,
	}
	
	if len(p.Extra) > 0 {
		copied.Extra = make(map[string]any, len(p.Extra))
		for k, v := range p.Extra {
			copied.Extra[k] = v
		}
	}
	
	return copied
}

func (p *Parameters) Merge(override *Parameters) {
	if override == nil {
		return
	}
	
	SetIfNotZero(&p.Temperature, override.Temperature)
	SetIfNotZero(&p.MaxTokens, override.MaxTokens)
	
	if len(override.Extra) > 0 {
		if p.Extra == nil {
			p.Extra = make(map[string]any)
		}
		for k, v := range override.Extra {
			p.Extra[k] = v
		}
	}
}

func (c *ChatCompletion) Copy() *ChatCompletion {
	if c == nil {
		return nil
	}

	copied := &ChatCompletion{
		APIFormat: c.APIFormat,
		Endpoint:  c.Endpoint,
		Context: Context{
			MaxInput:  c.Context.MaxInput,
			MaxOutput: c.Context.MaxOutput,
		},
		Features: Features{
			ToolUse:          c.Features.ToolUse,
			Thinking:         c.Features.Thinking,
			ThinkingLevels:   c.Features.ThinkingLevels,
			StructuredOutput: c.Features.StructuredOutput,
			AudioInput:       c.Features.AudioInput,
			ImageOutput:      c.Features.ImageOutput,
		},
		Parameters: c.Parameters.Copy(),
	}

	return copied
}

func (c *ChatCompletion) Merge(override *ChatCompletion) {
	SetIfNotZero(&c.APIFormat, override.APIFormat)
	SetIfNotZero(&c.Endpoint, override.Endpoint)

	SetIfNotZero(&c.Context.MaxInput, override.Context.MaxInput)
	SetIfNotZero(&c.Context.MaxOutput, override.Context.MaxOutput)

	SetIfNotZero(&c.Features.ToolUse, override.Features.ToolUse)
	SetIfNotZero(&c.Features.Thinking, override.Features.Thinking)
	SetIfNotZero(&c.Features.ThinkingLevels, override.Features.ThinkingLevels)
	SetIfNotZero(&c.Features.StructuredOutput, override.Features.StructuredOutput)
	SetIfNotZero(&c.Features.AudioInput, override.Features.AudioInput)
	SetIfNotZero(&c.Features.ImageOutput, override.Features.ImageOutput)

	c.Parameters.Merge(&override.Parameters)
}

type Context struct {
	MaxInput  int `yaml:"max_input" mapstructure:"max_input"`
	MaxOutput int `yaml:"max_output" mapstructure:"max_output"`
}

type Features struct {
	ToolUse          bool `yaml:"tool_use" mapstructure:"tool_use"`
	Thinking         bool `yaml:"thinking" mapstructure:"thinking"`
	ThinkingLevels   bool `yaml:"thinking_levels" mapstructure:"thinking_levels"`
	StructuredOutput bool `yaml:"structured_output" mapstructure:"structured_output"`
	AudioInput       bool `yaml:"audio_input" mapstructure:"audio_input"`
	ImageOutput      bool `yaml:"image_output" mapstructure:"image_output"`
}
