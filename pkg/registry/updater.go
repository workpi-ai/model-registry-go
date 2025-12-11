package registry

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	embed "github.com/workpi-ai/model-registry-go"
)

const (
	repoOwner       = "workpi-ai"
	repoName        = "model-registry"
	httpStatusOK    = 200
	defaultFilePerm = 0644
)

type Updater struct {
	configDir string
}

func NewUpdater(configDir string) *Updater {
	return &Updater{
		configDir: configDir,
	}
}

func (u *Updater) Update() error {
	latestVersion, err := u.getLatestVersion()
	if err != nil {
		return fmt.Errorf("get latest version: %w", err)
	}

	localVersion := u.getLocalVersion()

	if latestVersion == localVersion {
		u.updateCheckTime(localVersion)
		return nil
	}

	if err := u.downloadRelease(latestVersion); err != nil {
		return fmt.Errorf("download release: %w", err)
	}

	u.saveLocalVersion(latestVersion)

	return nil
}

func (u *Updater) getLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
		repoOwner, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func (u *Updater) getLocalVersion() string {
	path := filepath.Join(u.configDir, versionFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return embed.EmbedVersion
	}

	var v LocalVersion
	if err := json.Unmarshal(data, &v); err != nil {
		return embed.EmbedVersion
	}

	return v.Version
}

func (u *Updater) downloadRelease(version string) error {
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/providers.tar.gz",
		repoOwner, repoName, version)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}

	return u.extractTarGz(resp.Body)
}

func (u *Updater) extractTarGz(r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeDir {
			continue
		}

		if !strings.HasSuffix(header.Name, yamlExt) {
			continue
		}

		relPath := header.Name
		relPath = strings.TrimPrefix(relPath, providersDir+"/")
		relPath = strings.TrimPrefix(relPath, "./")

		destPath := filepath.Join(u.configDir, providersDir, relPath)

		if err = os.MkdirAll(filepath.Dir(destPath), defaultDirPerm); err != nil {
			return err
		}

		f, err := os.Create(destPath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(f, tr); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}

	return nil
}

func (u *Updater) saveLocalVersion(version string) error {
	v := LocalVersion{
		Version:     version,
		LastCheckAt: time.Now().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(u.configDir, versionFile)
	return os.WriteFile(path, data, defaultFilePerm)
}

func (u *Updater) updateCheckTime(version string) error {
	return u.saveLocalVersion(version)
}
