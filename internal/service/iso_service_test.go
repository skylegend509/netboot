package service

import (
	"mime/multipart"
	"netboot/internal/config"
	"os"
	"path/filepath"
	"testing"
)

func TestISOService(t *testing.T) {
	// Setup temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "netboot_iso_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		ISODir: tmpDir,
	}

	svc := NewISOService(cfg)

	// Test 1: Initial List should be empty
	if len(svc.List()) != 0 {
		t.Errorf("expected empty list, got %d", len(svc.List()))
	}

	// Test 2: Add Files manually to FS and Sync
	testFile := filepath.Join(tmpDir, "test_ubuntu.iso")
	if err := os.WriteFile(testFile, []byte("dummy iso content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Should not be visible yet if valid, but Sync calls are needed usually. 
	// However NewISOService calls Sync. We need to call Sync or trigger it.
	if err := svc.Sync(); err != nil {
		t.Fatal(err)
	}

	list := svc.List()
	if len(list) != 1 {
		t.Errorf("expected 1 file, got %d", len(list))
	} else if list[0].Name != "test_ubuntu.iso" {
		t.Errorf("expected 'test_ubuntu.iso', got '%s'", list[0].Name)
	}

	// Test 3: Delete
	if err := svc.Delete("test_ubuntu.iso"); err != nil {
		t.Fatal(err)
	}

	if len(svc.List()) != 0 {
		t.Errorf("expected empty list after delete, got %d", len(svc.List()))
	}
}

// MockFile implements multipart.File for testing
type MockFile struct {
	*os.File
}

func TestISOService_Add(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "netboot_iso_upload_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		ISODir: tmpDir,
	}
	svc := NewISOService(cfg)

	// Create a real temp file to simulate upload
	tmpFile, err := os.CreateTemp("", "upload_source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up source
	tmpFile.WriteString("iso content")
	tmpFile.Seek(0, 0)

	header := &multipart.FileHeader{
		Filename: "uploaded_arch.iso",
	}

	if err := svc.Add(tmpFile, header); err != nil {
		t.Fatal(err)
	}

	list := svc.List()
	if len(list) != 1 {
		t.Errorf("expected 1 file, got %d", len(list))
	}
	if list[0].Name != "uploaded_arch.iso" {
		t.Errorf("expected uploaded_arch.iso, got %s", list[0].Name)
	}
}
