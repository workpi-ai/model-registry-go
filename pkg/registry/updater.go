package registry

import (
	"strings"

	"github.com/workpi-ai/go-utils/ghrelease"
)

type providerFilter struct{}

func (f *providerFilter) ShouldExtract(filename string) bool {
	return strings.Contains(filename, providersDir+"/") && strings.HasSuffix(filename, yamlExt)
}

func NewUpdater(configDir string) *ghrelease.Updater {
	return ghrelease.NewUpdater(ghrelease.UpdaterConfig{
		RepoOwner:        repoOwner,
		RepoName:         repoName,
		DestDir:          configDir,
		MetadataFilename: versionFile,
		ExtractFilter:    &providerFilter{},
	})
}
