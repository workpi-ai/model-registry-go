package embed

import (
	"io/fs"
	
	registry "github.com/workpi-ai/model-registry"
)

var EmbedVersion = ""

func GetFS() (fs.FS, error) {
	return fs.Sub(registry.Providers, "providers")
}
