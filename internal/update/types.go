// Package update defines update-related types and helpers.
package update

import "time"

// Info contains information about a release.
type Info struct {
	Version         string              `json:"version"`
	ReleaseDate     time.Time           `json:"releaseDate"`
	ReleaseNotes    string              `json:"releaseNotes"`
	ReleaseNotesURL string              `json:"releaseNotesUrl,omitempty"`
	Mandatory       bool                `json:"mandatory"`
	Platforms       map[string]Platform `json:"platforms"`
}

// Platform contains platform-specific download information.
// key corresponds to "GOOS-GOARCH", e.g., "darwin-amd64", "windows-amd64".
type Platform struct {
	URL       string `json:"url"`
	Signature string `json:"signature"` // sha256:xxx
	Size      int64  `json:"size"`
}

// Status represents the current status of the update process.
type Status string

// Status values represent update lifecycle states.
const (
	StatusNone        Status = "none"
	StatusAvailable   Status = "available"
	StatusDownloading Status = "downloading"
	StatusReady       Status = "ready"
	StatusError       Status = "error"
)

// State reflects the current state of the update system,
// intended to be sent to the frontend.
type State struct {
	Status           Status  `json:"status"`
	CurrentVersion   string  `json:"currentVersion"`
	LatestVersion    string  `json:"latestVersion,omitempty"`
	UpdateInfo       *Info   `json:"updateInfo,omitempty"`
	DownloadProgress float64 `json:"downloadProgress,omitempty"` // 0-100
	Error            string  `json:"error,omitempty"`
}

// CheckOptions configures how updates are checked.
type CheckOptions struct {
	Prerelease bool `json:"prerelease"`
}
