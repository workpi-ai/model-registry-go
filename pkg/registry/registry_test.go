package registry

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()

	opts := Options{
		ConfigDir:     tmpDir,
		AutoUpdate:    false,
		CheckInterval: DefaultCheckInterval,
	}

	reg, err := New(opts)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if reg == nil {
		t.Fatal("expected non-nil registry")
	}

	if reg.Providers == nil {
		t.Fatal("expected non-nil Providers map")
	}
}

func TestNewWithInvalidConfigDir(t *testing.T) {
	opts := Options{
		ConfigDir: "",
	}

	_, err := New(opts)
	if err == nil {
		t.Fatal("expected error for empty ConfigDir")
	}
}

func TestCustomProviderMerge(t *testing.T) {
	tmpDir := t.TempDir()

	customProvider := &Provider{
		Name: "openai",
		Models: map[string]*Model{
			"gpt-4": {
				Name: "gpt-4",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  200000,
							MaxOutput: 20000,
						},
						Parameters: Parameters{
							Temperature: 0.7,
							MaxTokens:   4096,
						},
					},
				},
			},
		},
	}

	opts := Options{
		ConfigDir:     tmpDir,
		AutoUpdate:    false,
		CheckInterval: DefaultCheckInterval,
		Providers:     []*Provider{customProvider},
	}

	reg, err := New(opts)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	provider := reg.Provider("openai")
	if provider == nil {
		t.Fatal("expected to find openai provider")
	}

	model := provider.Models["gpt-4"]
	if model == nil {
		t.Fatal("expected to find gpt-4 model")
	}

	if model.APIs.ChatCompletion.Context.MaxInput != 200000 {
		t.Errorf("expected MaxInput=200000, got %d", model.APIs.ChatCompletion.Context.MaxInput)
	}
}

func TestCustomProviderNew(t *testing.T) {
	tmpDir := t.TempDir()

	customProvider := &Provider{
		Name: "custom-provider",
		Type: ProviderTypeAPI,
		Models: map[string]*Model{
			"custom-model": {
				Name: "custom-model",
				APIs: APIs{
					ChatCompletion: &ChatCompletion{
						Context: Context{
							MaxInput:  100000,
							MaxOutput: 10000,
						},
						Parameters: Parameters{
							Temperature: 0.5,
							MaxTokens:   2048,
						},
					},
				},
			},
		},
	}

	opts := Options{
		ConfigDir:     tmpDir,
		AutoUpdate:    false,
		CheckInterval: DefaultCheckInterval,
		Providers:     []*Provider{customProvider},
	}

	reg, err := New(opts)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	provider := reg.Provider("custom-provider")
	if provider == nil {
		t.Fatal("expected to find custom-provider")
	}

	if provider.Type != ProviderTypeAPI {
		t.Errorf("expected type=%s, got %s", ProviderTypeAPI, provider.Type)
	}

	model := provider.Models["custom-model"]
	if model == nil {
		t.Fatal("expected to find custom-model")
	}
}

func TestModelQuery(t *testing.T) {
	tmpDir := t.TempDir()

	customProvider := &Provider{
		Name: "test-provider",
		Models: map[string]*Model{
			"test-model": {
				Name: "test-model",
			},
		},
	}

	opts := Options{
		ConfigDir:     tmpDir,
		AutoUpdate:    false,
		CheckInterval: DefaultCheckInterval,
		Providers:     []*Provider{customProvider},
	}

	reg, err := New(opts)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	model := reg.Model("test-provider", "test-model")
	if model == nil {
		t.Fatal("expected to find test-model")
	}

	if model.Name != "test-model" {
		t.Errorf("expected model name=test-model, got %s", model.Name)
	}

	// Test non-existent provider
	model = reg.Model("non-existent", "test-model")
	if model != nil {
		t.Error("expected nil for non-existent provider")
	}

	// Test non-existent model
	model = reg.Model("test-provider", "non-existent")
	if model != nil {
		t.Error("expected nil for non-existent model")
	}
}

func TestMergeCustomProvidersEmptyName(t *testing.T) {
	tmpDir := t.TempDir()

	customProvider := &Provider{
		Name: "",
	}

	opts := Options{
		ConfigDir:     tmpDir,
		AutoUpdate:    false,
		CheckInterval: DefaultCheckInterval,
		Providers:     []*Provider{customProvider},
	}

	_, err := New(opts)
	if err == nil {
		t.Fatal("expected error for custom provider with empty name")
	}
}

func TestAutoUpdateLoop(t *testing.T) {
	tmpDir := t.TempDir()

	opts := Options{
		ConfigDir:     tmpDir,
		AutoUpdate:    true,
		CheckInterval: 100 * time.Millisecond,
	}

	reg, err := New(opts)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Let it run for a bit
	time.Sleep(150 * time.Millisecond)

	// Close should stop the goroutine
	if err := reg.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// Give it time to stop
	time.Sleep(50 * time.Millisecond)
}

func TestClose(t *testing.T) {
	tmpDir := t.TempDir()

	opts := Options{
		ConfigDir:     tmpDir,
		AutoUpdate:    true,
		CheckInterval: DefaultCheckInterval,
	}

	reg, err := New(opts)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	err = reg.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}
}
