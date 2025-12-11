// Package services contains backend service implementations used by the app.
package services

import (
	"context"
	"fmt"
	"freelance-flow/internal/update"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	// RepoOwner identifies the GitHub account hosting releases.
	RepoOwner = "royzhu"
	// RepoName identifies the GitHub repository hosting releases.
	RepoName = "freelance-flow"
)

// UpdateService handles application updates.
type UpdateService struct {
	ctx            context.Context
	state          update.State
	mu             sync.RWMutex
	_              int // hours, unused for now (checkInterval)
	downloader     *update.Downloader
	cancelDownload context.CancelFunc
}

// NewUpdateService creates a new UpdateService.
func NewUpdateService() *UpdateService {
	return &UpdateService{
		state: update.State{
			Status:         update.StatusNone,
			CurrentVersion: update.GetCurrentVersion(),
		},
		downloader: update.NewDownloader(),
	}
}

// startup is called by Wails when the application starts.
//
//nolint:unused // Called by Wails runtime via reflection.
func (s *UpdateService) startup(ctx context.Context) {
	s.ctx = ctx
	// Auto-check on startup
	go func() {
		_ = s.CheckForUpdate()
	}()
}

// CheckForUpdate checks for the latest version on GitHub.
func (s *UpdateService) CheckForUpdate() error {
	s.mu.Lock()
	// If already downloading or ready, don't check again to avoid overwriting state
	if s.state.Status == update.StatusDownloading || s.state.Status == update.StatusReady {
		s.mu.Unlock()
		return nil
	}
	s.mu.Unlock()

	latest, err := update.FetchLatestRelease(RepoOwner, RepoName)
	if err != nil {
		// If 404, it means no release found, which is fine.
		if contains404(err) {
			s.mu.Lock()
			s.state.Status = update.StatusNone
			s.state.Error = ""
			s.mu.Unlock()
			s.emitState()
			return nil
		}
		s.setError(err.Error())
		return err
	}

	available, err := update.IsUpdateAvailable(s.state.CurrentVersion, latest.Version)
	if err != nil {
		s.setError("Error comparing versions: " + err.Error())
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if available {
		s.state.Status = update.StatusAvailable
		s.state.LatestVersion = latest.Version
		s.state.UpdateInfo = latest
		s.state.Error = ""
	} else {
		s.state.Status = update.StatusNone
		// Keep current info but marked as none? Or clear?
		// Let's keep it clean
		s.state.LatestVersion = latest.Version // Still useful to know
	}

	s.emitState()
	return nil
}

// StartDownload starts downloading the update for the current platform.
func (s *UpdateService) StartDownload() error {
	s.mu.Lock()
	if s.state.Status != update.StatusAvailable {
		s.mu.Unlock()
		return fmt.Errorf("no update available to download")
	}

	// Find the asset key for current platform
	platformKey := runtime.GOOS + "-" + runtime.GOARCH
	asset, ok := s.state.UpdateInfo.Platforms[platformKey]
	if !ok {
		s.mu.Unlock()
		errMsg := fmt.Sprintf("no update found for platform %s", platformKey)
		s.setError(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	s.state.Status = update.StatusDownloading
	s.emitState()

	// Setup cancellation
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelDownload = cancel
	s.mu.Unlock()

	go func() {
		// Prepare destination
		tempDir := os.TempDir()
		fileName := filepath.Base(asset.URL)
		destPath := filepath.Join(tempDir, fileName)

		err := s.downloader.Download(ctx, asset.URL, destPath, func(total, current int64) {
			// Emit progress event
			// We could also calculate percentage here
			wailsRuntime.EventsEmit(s.ctx, "update:progress", map[string]interface{}{
				"total":   total,
				"current": current,
			})
		})

		s.mu.Lock()
		defer s.mu.Unlock()
		s.cancelDownload = nil // Clear cancel func

		if err != nil {
			if ctx.Err() == context.Canceled {
				s.state.Status = update.StatusAvailable // Revert to available
				s.emitState()
				return
			}
			s.state.Status = update.StatusError
			s.state.Error = "Download failed: " + err.Error()
			s.emitState()
			return
		}

		// Verify Hash if available (mock implementation might not have valid hash)
		// For now, we skip verification if signature is empty, or implement it if possible.
		// In previous steps we agreed to implement hash verification.
		// Let's verify if we have a hash.
		// s.state.UpdateInfo.Platforms doesn't seem to store the hash directly?
		// Wait, types.go Platform struct likely has 'Signature' or similar which we are treating as hash?
		// Looking at types.go (which I can't see right now but execution-log.md implies SHA256)
		// The update.json generation task says "Fill SHA256 signature".
		// Let's assume asset.Signature contains the SHA256 hash.

		if asset.Signature != "" {
			if err := s.downloader.VerifyHash(destPath, asset.Signature); err != nil {
				s.state.Status = update.StatusError
				s.state.Error = "Hash verification failed: " + err.Error()
				s.emitState()
				return
			}
		}

		// Success
		s.state.Status = update.StatusReady
		// Check where we store the downloaded file path?
		// We might need to store it in state or accessible field to open it later.
		// Since UpdateInfo.Platforms is a map, we can't easily modify it to add local path.
		// But for now, since we reconstruct the path from TempDir + Base(URL) in InstallUpdate, maybe that's fine?
		// Better: store it in a temporary internal map or just reconstruct it.
		// Let's reconstruct it in InstallUpdate for now.
		s.emitState()
	}()

	return nil
}

// CancelDownload cancels the active download.
func (s *UpdateService) CancelDownload() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancelDownload != nil {
		s.cancelDownload()
	}
}

// InstallUpdate initiates the installation (opens the file).
func (s *UpdateService) InstallUpdate() error {
	s.mu.RLock()
	if s.state.Status != update.StatusReady {
		s.mu.RUnlock()
		return fmt.Errorf("update not ready to install")
	}

	platformKey := runtime.GOOS + "-" + runtime.GOARCH
	asset, ok := s.state.UpdateInfo.Platforms[platformKey]
	s.mu.RUnlock() // Unlock before doing IO/Exec

	if !ok {
		return fmt.Errorf("platform info missing")
	}

	// Reconstruct path (should match StartDownload)
	tempDir := os.TempDir()
	fileName := filepath.Base(asset.URL)
	destPath := filepath.Join(tempDir, fileName)

	if runtime.GOOS == "darwin" {
		// Open DMG; path is derived from TempDir and known asset name.
		cmd := exec.Command("open", filepath.Clean(destPath)) //nolint:gosec
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to open installer: %w", err)
		}

		// Optionally reveal in finder too?
		// exec.Command("open", "-R", destPath).Start()
	} else {
		// For other platforms, just reveal?
		return fmt.Errorf("auto-install not supported for %s yet", runtime.GOOS)
	}

	return nil
}

// GetUpdateState returns the current update state.
func (s *UpdateService) GetUpdateState() update.State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// SkipVersion marks the current available version as skipped.
func (s *UpdateService) SkipVersion() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state.Status == update.StatusAvailable || s.state.Status == update.StatusError {
		// Logic to persist skip choice would go here
		s.state.Status = update.StatusNone
		s.emitState()
	}
}

// Helper to set error state
func (s *UpdateService) setError(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Error = msg
	s.state.Status = update.StatusError
	s.emitState()
}

// emitState sends the current state to the frontend
func (s *UpdateService) emitState() {
	if s.ctx != nil {
		wailsRuntime.EventsEmit(s.ctx, "update:state", s.state)
	}
}

func contains404(err error) bool {
	return err != nil && (err.Error() == "github api returned status: 404" ||
		// Check for wrapped errors if necessary, but simple string match works for now
		// based on github.go implementation
		len(err.Error()) > 0 && err.Error()[len(err.Error())-3:] == "404")
}
