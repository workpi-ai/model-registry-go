package registry

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestModelValidate(t *testing.T) {
	tests := []struct {
		name      string
		model     *Model
		wantError bool
	}{
		{
			name: "valid model",
			model: &Model{
				Name: "gpt-4",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  8000,
							MaxOutput: 4000,
						},
						Parameters: Parameters{
							Temperature: 0.7,
							MaxTokens:   1000,
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "empty name",
			model: &Model{
				Name: "",
			},
			wantError: true,
		},
		{
			name: "negative max_input",
			model: &Model{
				Name: "test-model",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  -100,
							MaxOutput: 1000,
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "zero max_input",
			model: &Model{
				Name: "test-model",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  0,
							MaxOutput: 1000,
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "negative max_output",
			model: &Model{
				Name: "test-model",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  1000,
							MaxOutput: -100,
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "zero max_output",
			model: &Model{
				Name: "test-model",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  1000,
							MaxOutput: 0,
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "max_output exceeds max_input",
			model: &Model{
				Name: "test-model",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  1000,
							MaxOutput: 2000,
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "model without chat completion API",
			model: &Model{
				Name: "test-model",
				APIs: APIs{},
			},
			wantError: false,
		},
		{
			name: "negative max_tokens",
			model: &Model{
				Name: "test-model",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  1000,
							MaxOutput: 500,
						},
						Parameters: Parameters{
							MaxTokens: -100,
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "zero max_tokens",
			model: &Model{
				Name: "test-model",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  1000,
							MaxOutput: 500,
						},
						Parameters: Parameters{
							MaxTokens: 0,
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "positive max_tokens",
			model: &Model{
				Name: "test-model",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  1000,
							MaxOutput: 500,
						},
						Parameters: Parameters{
							MaxTokens: 100,
						},
					},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			if tt.wantError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestLegacyReasoningLevelYAMLDoesNotPopulateReasoningEffort(t *testing.T) {
	var chatCompletion ChatCompletion
	if err := yaml.Unmarshal([]byte(`
features:
  reasoning_levels:
    - high
parameters:
  reasoning_level: high
`), &chatCompletion); err != nil {
		t.Fatalf("unmarshal legacy chat completion: %v", err)
	}

	if chatCompletion.Parameters.ReasoningEffort != "" {
		t.Fatalf("legacy reasoning_level populated canonical field: %q", chatCompletion.Parameters.ReasoningEffort)
	}
	if len(chatCompletion.Features.ReasoningEfforts) != 0 {
		t.Fatalf("legacy reasoning_levels populated canonical field: %v", chatCompletion.Features.ReasoningEfforts)
	}
}

func TestChatCompletionReasoningEffortYAML(t *testing.T) {
	var chatCompletion ChatCompletion
	if err := yaml.Unmarshal([]byte(`
features:
  reasoning: true
  reasoning_efforts:
    - low
    - high
parameters:
  reasoning_effort: high
`), &chatCompletion); err != nil {
		t.Fatalf("unmarshal chat completion: %v", err)
	}

	if chatCompletion.Parameters.ReasoningEffort != "high" {
		t.Fatalf("expected reasoning effort high, got %q", chatCompletion.Parameters.ReasoningEffort)
	}
	if !reflect.DeepEqual(chatCompletion.Features.ReasoningEfforts, []string{"low", "high"}) {
		t.Fatalf("unexpected reasoning efforts: %v", chatCompletion.Features.ReasoningEfforts)
	}

	copied := chatCompletion.Copy()
	copied.Features.ReasoningEfforts[0] = "medium"
	if chatCompletion.Features.ReasoningEfforts[0] != "low" {
		t.Fatalf("copy shares reasoning efforts with original: %v", chatCompletion.Features.ReasoningEfforts)
	}
	if copied.Parameters.ReasoningEffort != "high" {
		t.Fatalf("copy lost reasoning effort: %q", copied.Parameters.ReasoningEffort)
	}

	chatCompletion.Merge(&ChatCompletion{
		Features: Features{ReasoningEfforts: []string{"max"}},
		Parameters: Parameters{
			ReasoningEffort: "max",
		},
	})
	if !reflect.DeepEqual(chatCompletion.Features.ReasoningEfforts, []string{"max"}) {
		t.Fatalf("merge did not replace reasoning efforts: %v", chatCompletion.Features.ReasoningEfforts)
	}
	if chatCompletion.Parameters.ReasoningEffort != "max" {
		t.Fatalf("merge did not replace reasoning effort: %q", chatCompletion.Parameters.ReasoningEffort)
	}
}
