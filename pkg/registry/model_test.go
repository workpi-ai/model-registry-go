package registry

import (
	"testing"
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