package registry

import (
	"strings"

	"github.com/workpi-ai/go-utils/ghrelease"
)

type providerPathTransformer struct {
	destDir string
}

func (t *providerPathTransformer) Transform(filename string) string {
	if strings.Contains(filename, providersDir+"/") && strings.HasSuffix(filename, yamlExt) {
		return filename
	}
	return ""
}

func NewUpdater(destDir string) (*ghrelease.Updater, error) {
	return ghrelease.NewUpdater(ghrelease.UpdaterConfig{
		RepoOwner:    repoOwner,
		RepoName:     repoName,
		MetadataFile: destDir + "/" + versionFile,
		Targets: []ghrelease.ExtractTarget{
			{
				PathTransformer: &providerPathTransformer{destDir: destDir},
				DestDir:         destDir,
			},
		},
	})
}
