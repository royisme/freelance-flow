// Package main provides a tool for automating the release process.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/blang/semver/v4"
)

type WailsConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	// Other fields ignored
}

func main() {
	var part string
	flag.StringVar(&part, "part", "patch", "Part of version to bump (patch, minor, major)")
	flag.Parse()

	// 1. Read current version from wails.json
	wailsContent, err := os.ReadFile("wails.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to read wails.json: %v", err))
	}

	// We decode into a map to preserve other fields structure if we were rewriting the whole file
	// But JSON marshalling might reorder keys.
	// To minimize diff, let's use a struct just for reading version, and string replacement for writing?
	// Or better: Use semver to increment, then regex replace in file content.
	// This preserves comments and formatting better than re-encoding.

	var config WailsConfig
	if err := json.Unmarshal(wailsContent, &config); err != nil {
		panic(err)
	}

	currentVer, err := semver.Parse(config.Version)
	if err != nil {
		panic(err)
	}

	// 2. Bump version
	newVer := currentVer
	switch part {
	case "major":
		newVer.Major++
		newVer.Minor = 0
		newVer.Patch = 0
	case "minor":
		newVer.Minor++
		newVer.Patch = 0
	case "patch":
		newVer.Patch++
	default:
		panic("Invalid bump part: " + part)
	}

	fmt.Printf("Bumping version: %s -> %s\n", currentVer, newVer)

	// 3. Update wails.json
	// Simple string replacement to avoid reformatting entire JSON
	// Assuming "version": "1.0.0" format
	newContent := strings.Replace(string(wailsContent),
		fmt.Sprintf(`"version": "%s"`, currentVer),
		fmt.Sprintf(`"version": "%s"`, newVer),
		1)

	// #nosec G703 -- wails.json is a hardcoded path, not user input
	if err := os.WriteFile("wails.json", []byte(newContent), 0600); err != nil {
		panic(err)
	}

	// 4. Git operations
	runCommand("git", "add", "wails.json")

	commitMsg := fmt.Sprintf("Bump version to v%s", newVer)
	runCommand("git", "commit", "-m", commitMsg)

	tagName := fmt.Sprintf("v%s", newVer)
	runCommand("git", "tag", tagName)

	fmt.Printf("\nSuccess! \nRun 'git push origin main %s' to trigger release.\n", tagName)
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("> %s %s\n", name, strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
