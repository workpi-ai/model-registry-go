# Model Registry Go SDK

Go SDK for [model-registry](https://github.com/workpi-ai/model-registry) - A centralized metadata repository for AI models.

## Features

- 🚀 **Embedded data**: Works offline with embedded registry
- 🔄 **Auto-update**: Automatically checks for updates on startup
- 📦 **Lightweight**: Only ~150KB of YAML data
- 🌐 **Multi-provider**: Supports OpenAI, Anthropic, Google, DeepSeek, and more

## Installation

```bash
go get github.com/workpi-ai/model-registry-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/workpi-ai/model-registry-go/pkg/registry"
)

func main() {
    home, _ := os.UserHomeDir()
    configDir := filepath.Join(home, ".codev", "configs")
    
    reg, err := registry.New(registry.Options{
        ConfigDir:  configDir,
        AutoUpdate: false,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Get model info
    model, _ := reg.GetModel("openai/gpt-4o")
    fmt.Printf("Max Input: %d\n", model.APIs["chat_completion"].Context.MaxInput)
    
    // List providers
    providers := reg.ListProviders()
    fmt.Printf("Total providers: %d\n", len(providers))
}
```

## Configuration

```go
import (
    "os"
    "path/filepath"
    "time"
)

home, _ := os.UserHomeDir()
configDir := filepath.Join(home, ".codev", "configs")

reg, err := registry.New(registry.Options{
    ConfigDir:     configDir,        // Required: Local cache directory
    AutoUpdate:    true,              // Auto-check for updates
    CheckInterval: 1 * time.Hour,     // Check interval (default: 1h)
})
```

## Data Priority

1. **Local cache** (`$HOME/.codev/configs/providers/`) - Downloaded from GitHub Release
2. **Embedded data** - Bundled from the model-registry Go module dependency

## API Reference

### Get Provider

```go
provider, err := reg.GetProvider("openai")
fmt.Printf("Base URL: %s\n", provider.BaseURL)
fmt.Printf("Auth Type: %s\n", provider.AuthType)
```

### Get Model

```go
model, err := reg.GetModel("openai/gpt-4o")
fmt.Printf("Model: %s\n", model.Name)
fmt.Printf("Provider: %s\n", model.Provider)
fmt.Printf("Max Input: %d\n", model.APIs["chat_completion"].Context.MaxInput)
fmt.Printf("Max Output: %d\n", model.APIs["chat_completion"].Context.MaxOutput)
fmt.Printf("Tool Use: %v\n", model.APIs["chat_completion"].Features.ToolUse)
```

### List Providers

```go
providers := reg.ListProviders()
for _, name := range providers {
    fmt.Println(name)
}
```

### List Models

```go
// List all models
allModels := reg.ListModels("")

// List models for specific provider
openaiModels := reg.ListModels("openai")
```

### Force Update

```go
// Manually trigger update
err := reg.ForceUpdate()
if err != nil {
    log.Fatal(err)
}
```

## Development

### Clone Repository

```bash
git clone https://github.com/workpi-ai/model-registry-go.git
```

### Build

```bash
make build
```

### Test

```bash
make test
```

## How It Works

1. **Compile Time**: Embeds registry data from the model-registry Go module
2. **Runtime**: 
   - On startup, checks for updates from GitHub Release (if AutoUpdate is enabled)
   - Downloads new version to `$HOME/.codev/configs/providers/` if available
   - Loads data with priority: local cache > embedded data

## License

MIT
