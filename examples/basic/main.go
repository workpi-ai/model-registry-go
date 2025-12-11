package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/workpi-ai/model-registry-go/pkg/registry"
)

func main() {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".codev", "configs")

	reg, err := registry.New(registry.Options{
		ConfigDir:     configDir,
		AutoUpdate:    true,
		CheckInterval: 1 * time.Hour,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer reg.Close()

	// List all providers
	providers := reg.ProviderList()
	fmt.Printf("Available providers (%d):\n", len(providers))
	for _, p := range providers {
		fmt.Printf("  - %s: %s\n", p.Name, p.Description)
	}
	fmt.Println()

	// Get specific model
	model := reg.Model("openai", "gpt-4o")
	if model == nil {
		log.Fatal("model not found")
	}

	fmt.Printf("Model: %s\n", model.Name)
	fmt.Printf("Provider: %s\n", model.Provider.Name)
	fmt.Printf("Max Input: %d\n", model.APIs.ChatCompletion.Context.MaxInput)
	fmt.Printf("Max Output: %d\n", model.APIs.ChatCompletion.Context.MaxOutput)
}

