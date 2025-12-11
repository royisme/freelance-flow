// Package main generates update.json metadata for releases.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"freelance-flow/internal/update"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var (
		version      string
		releaseNotes string
		repoOwner    string
		repoName     string
		output       string
	)

	flag.StringVar(&version, "version", "", "Release version (e.g. 1.0.0)")
	flag.StringVar(&releaseNotes, "notes", "", "Release notes")
	flag.StringVar(&repoOwner, "owner", "royzhu", "Repository owner")
	flag.StringVar(&repoName, "repo", "freelance-flow", "Repository name")
	flag.StringVar(&output, "output", "update.json", "Output file path")
	flag.Parse()

	if version == "" {
		panic("version is required")
	}

	// Clean version
	cleanVersion := strings.TrimPrefix(version, "v")

	// Parse arguments for files
	platforms := make(map[string]update.Platform)

	for _, file := range flag.Args() {
		// Expect file name format: FreelanceFlow-{version}-{os}-{arch}.ext
		// But in CI we might just list specific files.
		// Let's rely on simple mapping based on filename content.

		info, err := os.Stat(file)
		if err != nil {
			fmt.Printf("Warning: skipping file %s: %v\n", file, err)
			continue
		}

		hash, err := calculateSHA256(file)
		if err != nil {
			panic(fmt.Sprintf("Failed to calculate hash for %s: %v", file, err))
		}

		filename := filepath.Base(file)
		url := fmt.Sprintf("https://github.com/%s/%s/releases/download/v%s/%s",
			repoOwner, repoName, cleanVersion, filename)

		platform := update.Platform{
			URL:       url,
			Signature: "sha256:" + hash,
			Size:      info.Size(),
		}

		if strings.Contains(filename, "darwin") {
			if strings.Contains(filename, "amd64") {
				platforms["darwin-amd64"] = platform
			} else if strings.Contains(filename, "arm64") {
				platforms["darwin-arm64"] = platform
			}
		} else if strings.Contains(filename, "windows") {
			if strings.Contains(filename, "amd64") {
				platforms["windows-amd64"] = platform
			}
		}
	}

	info := update.Info{
		Version:         cleanVersion,
		ReleaseDate:     time.Now().UTC(),
		ReleaseNotes:    releaseNotes,
		ReleaseNotesURL: fmt.Sprintf("https://github.com/%s/%s/releases/tag/v%s", repoOwner, repoName, cleanVersion),
		Platforms:       platforms,
	}

	// G304: Risk accepted for CLI tool input. G306: Use 0600.
	file, err := os.OpenFile(filepath.Clean(output), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(info); err != nil {
		panic(err)
	}

	fmt.Printf("Generated %s for version %s\n", output, cleanVersion)
}

func calculateSHA256(path string) (string, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
