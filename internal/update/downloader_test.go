package update

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownload(t *testing.T) {
	content := "Hello, World! This is a test file for downloading."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		if _, err := fmt.Fprint(w, content); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	d := NewDownloader()
	tmpDir := t.TempDir()
	destPath := filepath.Join(tmpDir, "downloaded.txt")

	progressCalled := false
	onProgress := func(total, current int64) {
		progressCalled = true
		if total != int64(len(content)) {
			t.Errorf("Expected total %d, got %d", len(content), total)
		}
		_ = current // ensure parameters are exercised
	}

	err := d.Download(context.Background(), server.URL, destPath, onProgress)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	if !progressCalled {
		t.Error("onProgress was not called")
	}

	downloadedInfo, err := os.ReadFile(filepath.Clean(destPath))
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	if string(downloadedInfo) != content {
		t.Errorf("Content mismatch. Expected %q, got %q", content, string(downloadedInfo))
	}
}

func TestVerifyHash(t *testing.T) {
	content := "Hash test content"
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "hash_test.txt")
	err := os.WriteFile(filePath, []byte(content), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hasher := sha256.New()
	hasher.Write([]byte(content))
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	d := NewDownloader()

	// Test valid hash
	if err := d.VerifyHash(filePath, expectedHash); err != nil {
		t.Errorf("VerifyHash failed for valid hash: %v", err)
	}

	// Test invalid hash
	if err := d.VerifyHash(filePath, "badhash"); err == nil {
		t.Error("VerifyHash should have failed for invalid hash")
	}
}
