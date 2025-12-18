package registry

import (
	"testing"
)

func TestProviderMerge(t *testing.T) {
	base := &Provider{
		Name:        "openai",
		Type:        ProviderTypeAPI,
		AuthType:    AuthTypeAPIKey,
		APIKey:      "base-key",
		BaseURL:     "https://api.openai.com",
		Description: "Base description",
		Models:      make(map[string]*Model),
	}

	base.Models["gpt-4"] = &Model{
		Name:     "gpt-4",
		Provider: base,
		APIs: APIs{
			ChatCompletion: &ChatCompletion{
				APIFormat: APIFormatOpenAI,
				Endpoint:  "/v1/chat/completions",
				Context: Context{
					MaxInput:  8000,
					MaxOutput: 4000,
				},
				Features: Features{
					ToolUse:          true,
					Thinking:         false,
					StructuredOutput: true,
				},
				Parameters: Parameters{
					Temperature: 0.7,
				},
			},
		},
	}

	override := &Provider{
		Name:        "openai",
		APIKey:      "override-key",
		BaseURL:     "https://custom.api.com",
		Description: "",
		Models:      make(map[string]*Model),
	}

	override.Models["gpt-4"] = &Model{
		Name: "gpt-4",
		APIs: APIs{
			ChatCompletion: &ChatCompletion{
				Context: Context{
					MaxOutput: 8000,
				},
				Parameters: Parameters{
					MaxTokens: 1000,
				},
			},
		},
	}

	override.Models["gpt-3.5"] = &Model{
		Name: "gpt-3.5",
		APIs: APIs{
			ChatCompletion: &ChatCompletion{
				APIFormat: APIFormatOpenAI,
				Endpoint:  "/v1/chat/completions",
			},
		},
	}

	base.Merge(override)

	if base.APIKey != "override-key" {
		t.Errorf("Expected APIKey to be 'override-key', got '%s'", base.APIKey)
	}

	if base.BaseURL != "https://custom.api.com" {
		t.Errorf("Expected BaseURL to be 'https://custom.api.com', got '%s'", base.BaseURL)
	}

	if base.Description != "Base description" {
		t.Errorf("Expected Description to remain 'Base description', got '%s'", base.Description)
	}

	if len(base.Models) != 2 {
		t.Errorf("Expected 2 models, got %d", len(base.Models))
	}

	gpt4 := base.Models["gpt-4"]
	if gpt4 == nil {
		t.Fatal("Expected gpt-4 model to exist")
	}

	if gpt4.APIs.ChatCompletion.Context.MaxInput != 8000 {
		t.Errorf("Expected MaxInput to remain 8000, got %d", gpt4.APIs.ChatCompletion.Context.MaxInput)
	}

	if gpt4.APIs.ChatCompletion.Context.MaxOutput != 8000 {
		t.Errorf("Expected MaxOutput to be merged to 8000, got %d", gpt4.APIs.ChatCompletion.Context.MaxOutput)
	}

	if gpt4.APIs.ChatCompletion.Parameters.Temperature != 0.7 {
		t.Errorf("Expected temperature to remain 0.7, got %v", gpt4.APIs.ChatCompletion.Parameters.Temperature)
	}

	if gpt4.APIs.ChatCompletion.Parameters.MaxTokens != 1000 {
		t.Errorf("Expected max_tokens to be merged to 1000, got %v", gpt4.APIs.ChatCompletion.Parameters.MaxTokens)
	}

	gpt35 := base.Models["gpt-3.5"]
	if gpt35 == nil {
		t.Fatal("Expected gpt-3.5 model to exist")
	}

	if gpt35.Provider != base {
		t.Error("Expected gpt-3.5 Provider to be set to base")
	}
}

func TestModelMerge(t *testing.T) {
	base := &Model{
		Name: "gpt-4",
		APIs: APIs{
			ChatCompletion: &ChatCompletion{
				APIFormat: APIFormatOpenAI,
				Endpoint:  "/v1/chat/completions",
				Context: Context{
					MaxInput:  8000,
					MaxOutput: 4000,
				},
				Features: Features{
					ToolUse:          true,
					Thinking:        false,
					StructuredOutput: true,
				},
			},
		},
	}

	override := &Model{
		Name: "gpt-4-turbo",
		APIs: APIs{
			ChatCompletion: &ChatCompletion{
				Context: Context{
					MaxOutput: 8000,
				},
				Features: Features{
					Thinking: true,
				},
			},
		},
	}

	base.Merge(override)

	if base.Name != "gpt-4-turbo" {
		t.Errorf("Expected Name to be 'gpt-4-turbo', got '%s'", base.Name)
	}

	if base.APIs.ChatCompletion.Context.MaxInput != 8000 {
		t.Errorf("Expected MaxInput to remain 8000, got %d", base.APIs.ChatCompletion.Context.MaxInput)
	}

	if base.APIs.ChatCompletion.Context.MaxOutput != 8000 {
		t.Errorf("Expected MaxOutput to be merged to 8000, got %d", base.APIs.ChatCompletion.Context.MaxOutput)
	}

	if !base.APIs.ChatCompletion.Features.Thinking {
		t.Error("Expected Reasoning to be merged to true")
	}

	if !base.APIs.ChatCompletion.Features.ToolUse {
		t.Error("Expected ToolUse to remain true")
	}
}

func TestChatCompletionMerge(t *testing.T) {
	topP09 := 0.9
	base := &ChatCompletion{
		APIFormat: APIFormatOpenAI,
		Endpoint:  "/v1/chat/completions",
		Context: Context{
			MaxInput:  8000,
			MaxOutput: 4000,
		},
		Features: Features{
			ToolUse:          true,
			Thinking:        false,
			StructuredOutput: true,
		},
		Parameters: Parameters{
			Temperature: 0.7,
			Extra: map[string]any{
				"top_p": topP09,
			},
		},
	}

	override := &ChatCompletion{
		Endpoint: "/v1/chat/completions/new",
		Context: Context{
			MaxOutput: 8000,
		},
		Features: Features{
			Thinking: true,
		},
		Parameters: Parameters{
			Temperature: 0.5,
			MaxTokens:   1000,
		},
	}

	base.Merge(override)

	if base.Endpoint != "/v1/chat/completions/new" {
		t.Errorf("Expected Endpoint to be merged, got '%s'", base.Endpoint)
	}

	if base.Context.MaxInput != 8000 {
		t.Errorf("Expected MaxInput to remain 8000, got %d", base.Context.MaxInput)
	}

	if base.Context.MaxOutput != 8000 {
		t.Errorf("Expected MaxOutput to be merged to 8000, got %d", base.Context.MaxOutput)
	}

	if !base.Features.Thinking {
		t.Error("Expected Reasoning to be merged to true")
	}

	if !base.Features.ToolUse {
		t.Error("Expected ToolUse to remain true")
	}

	if base.Parameters.Temperature != 0.5 {
		t.Errorf("Expected temperature to be merged to 0.5, got %v", base.Parameters.Temperature)
	}

	if base.Parameters.Extra["top_p"] != 0.9 {
		t.Errorf("Expected top_p to remain 0.9, got %v", base.Parameters.Extra["top_p"])
	}

	if base.Parameters.MaxTokens != 1000 {
		t.Errorf("Expected max_tokens to be merged to 1000, got %v", base.Parameters.MaxTokens)
	}
}

func TestProviderCopy(t *testing.T) {
	original := &Provider{
		Name:        "openai",
		Type:        ProviderTypeAPI,
		AuthType:    AuthTypeAPIKey,
		APIKey:      "test-key",
		BaseURL:     "https://api.openai.com",
		Description: "Test provider",
		Models:      make(map[string]*Model),
	}

	original.Models["gpt-4"] = &Model{
		Name:     "gpt-4",
		Provider: original,
		APIs: APIs{
			ChatCompletion: &ChatCompletion{
				APIFormat: APIFormatOpenAI,
				Endpoint:  "/v1/chat/completions",
			},
		},
	}

	copied := original.Copy()

	if copied == original {
		t.Error("Copy should return a new instance")
	}

	if copied.Name != original.Name {
		t.Error("Copied Name should match original")
	}

	if copied.Models["gpt-4"] == original.Models["gpt-4"] {
		t.Error("Copied model should be a new instance")
	}

	if copied.Models["gpt-4"].Provider != copied {
		t.Error("Copied model's Provider should point to copied provider")
	}

	copied.APIKey = "new-key"
	if original.APIKey == "new-key" {
		t.Error("Modifying copy should not affect original")
	}
}

func TestProviderValidate(t *testing.T) {
	tests := []struct {
		name      string
		provider  *Provider
		wantError bool
	}{
		{
			name: "valid provider with api_key",
			provider: &Provider{
				Name:     "test",
				AuthType: AuthTypeAPIKey,
				APIKey:   "sk-xxx",
			},
			wantError: false,
		},
		{
			name: "valid provider without api_key auth",
			provider: &Provider{
				Name:     "test",
				AuthType: AuthTypeOAuth2,
			},
			wantError: false,
		},
		{
			name: "empty name",
			provider: &Provider{
				Name: "",
			},
			wantError: true,
		},
		{
			name: "api_key auth without key",
			provider: &Provider{
				Name:     "test",
				AuthType: AuthTypeAPIKey,
				APIKey:   "",
			},
			wantError: true,
		},
		{
			name: "negative max_input",
			provider: &Provider{
				Name: "test",
				Models: map[string]*Model{
					"model1": {
						Name: "model1",
						APIs: APIs{
							ChatCompletion: &ChatCompletion{
								Context: Context{
									MaxInput:  -100,
									MaxOutput: 1000,
								},
							},
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "negative max_output",
			provider: &Provider{
				Name: "test",
				Models: map[string]*Model{
					"model1": {
						Name: "model1",
						APIs: APIs{
							ChatCompletion: &ChatCompletion{
								Context: Context{
									MaxInput:  1000,
									MaxOutput: -100,
								},
							},
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "max_output exceeds max_input",
			provider: &Provider{
				Name: "test",
				Models: map[string]*Model{
					"model1": {
						Name: "model1",
						APIs: APIs{
							ChatCompletion: &ChatCompletion{
								Context: Context{
									MaxInput:  1000,
									MaxOutput: 2000,
								},
							},
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "valid context values",
			provider: &Provider{
				Name: "test",
				Models: map[string]*Model{
					"model1": {
						Name: "model1",
						APIs: APIs{
							ChatCompletion: &ChatCompletion{
								Context: Context{
									MaxInput:  2000,
									MaxOutput: 1000,
								},
								Parameters: Parameters{
									Temperature: 0.7,
									MaxTokens:   1000,
								},
							},
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "zero max_input",
			provider: &Provider{
				Name: "test",
				Models: map[string]*Model{
					"model1": {
						Name: "model1",
						APIs: APIs{
							ChatCompletion: &ChatCompletion{
								Context: Context{
									MaxInput:  0,
									MaxOutput: 1000,
								},
							},
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "zero max_output",
			provider: &Provider{
				Name: "test",
				Models: map[string]*Model{
					"model1": {
						Name: "model1",
						APIs: APIs{
							ChatCompletion: &ChatCompletion{
								Context: Context{
									MaxInput:  1000,
									MaxOutput: 0,
								},
							},
						},
					},
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.provider.Validate()
			if tt.wantError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}
