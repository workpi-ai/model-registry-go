package registry

import (
	"testing"
)

func TestSetIfNotZero(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		source   string
		expected string
	}{
		{
			name:     "overwrite with non-empty string",
			initial:  "old",
			source:   "new",
			expected: "new",
		},
		{
			name:     "skip empty string",
			initial:  "old",
			source:   "",
			expected: "old",
		},
		{
			name:     "overwrite empty with non-empty",
			initial:  "",
			source:   "new",
			expected: "new",
		},
		{
			name:     "both empty",
			initial:  "",
			source:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.initial
			SetIfNotZero(&target, tt.source)
			if target != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, target)
			}
		})
	}
}

func TestSetIfNotZeroWithBool(t *testing.T) {
	tests := []struct {
		name     string
		initial  bool
		source   bool
		expected bool
	}{
		{
			name:     "overwrite false with true",
			initial:  false,
			source:   true,
			expected: true,
		},
		{
			name:     "skip false source",
			initial:  true,
			source:   false,
			expected: true,
		},
		{
			name:     "overwrite true with true",
			initial:  true,
			source:   true,
			expected: true,
		},
		{
			name:     "both false",
			initial:  false,
			source:   false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.initial
			SetIfNotZero(&target, tt.source)
			if target != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, target)
			}
		})
	}
}

func TestSetIfNotZeroWithCustomTypes(t *testing.T) {
	// Test with ProviderType
	var providerType ProviderType = "api"
	SetIfNotZero(&providerType, ProviderType("subscription"))
	if providerType != "subscription" {
		t.Errorf("expected subscription, got %s", providerType)
	}

	// Test with empty ProviderType
	providerType = "api"
	SetIfNotZero(&providerType, ProviderType(""))
	if providerType != "api" {
		t.Errorf("expected api, got %s", providerType)
	}

	// Test with AuthType
	var authType AuthType = "oauth2"
	SetIfNotZero(&authType, AuthType("api_key"))
	if authType != "api_key" {
		t.Errorf("expected api_key, got %s", authType)
	}
}

func TestSetIfNotZeroWithInt(t *testing.T) {
	tests := []struct {
		name     string
		initial  int
		source   int
		expected int
	}{
		{
			name:     "overwrite with positive number",
			initial:  100,
			source:   200,
			expected: 200,
		},
		{
			name:     "skip zero",
			initial:  100,
			source:   0,
			expected: 100,
		},
		{
			name:     "overwrite with negative number",
			initial:  100,
			source:   -50,
			expected: -50,
		},
		{
			name:     "overwrite zero with positive",
			initial:  0,
			source:   100,
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.initial
			SetIfNotZero(&target, tt.source)
			if target != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, target)
			}
		})
	}
}