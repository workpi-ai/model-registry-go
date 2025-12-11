package registry

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v68/github"
	embed "github.com/workpi-ai/model-registry-go"
)

const (
	repoOwner       = "workpi-ai"
	repoName        = "model-registry"
	requestTimeout  = 3 * time.Second
	downloadTimeout = 30 * time.Second
)

type Updater struct {
	configDir string
	client    *github.Client
}

func NewUpdater(configDir string) *Updater {
	return &Updater{
		configDir: configDir,
		client:    github.NewClient(nil),
	}
}

func (u *Updater) Update() error {
	latestVersion, err := u.getLatestVersion()
	if err != nil {
		return fmt.Errorf("get latest version: %w", err)
	}

	localVersion := u.getLocalVersion()

	if latestVersion == localVersion {
		u.saveLocalVersion(localVersion)
		return nil
	}

	if err := u.downloadRelease(latestVersion); err != nil {
		return fmt.Errorf("download release: %w", err)
	}

	u.saveLocalVersion(latestVersion)

	return nil
}

func (u *Updater) getLatestVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	
	release, _, err := u.client.Repositories.GetLatestRelease(ctx, repoOwner, repoName)
	if err != nil {
		return "", fmt.Errorf("failed to get latest release: %w", err)
	}

	if release.TagName == nil {
		return "", fmt.Errorf("release tag name is nil")
	}

	return *release.TagName, nil
}

func (u *Updater) getLocalVersion() string {
	path := filepath.Join(u.configDir, versionFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return embed.EmbedVersion
	}

	var v Metadata
	if err := json.Unmarshal(data, &v); err != nil {
		return embed.EmbedVersion
	}

	return v.Version
}

func (u *Updater) downloadRelease(version string) error {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()
	
	release, _, err := u.client.Repositories.GetReleaseByTag(ctx, repoOwner, repoName, version)
	if err != nil {
		return fmt.Errorf("failed to get release info: %w", err)
	}

	if release.ZipballURL == nil {
		return fmt.Errorf("release zipball_url is nil")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, *release.ZipballURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download zipball: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, *release.ZipballURL)
	}

	return u.extractZip(resp.Body)
}

func (u *Updater) extractZip(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name, yamlExt) {
			continue
		}

		idx := strings.Index(file.Name, providersDir+"/")
		if idx == -1 {
			continue
		}

		relPath := file.Name[idx+len(providersDir)+1:]
		destPath := filepath.Join(u.configDir, providersDir, relPath)

		if err = os.MkdirAll(filepath.Dir(destPath), defaultDirPerm); err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		f, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(f, rc); err != nil {
			return err
		}
	}

	return nil
}

func (u *Updater) saveLocalVersion(version string) error {
	v := Metadata{
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
