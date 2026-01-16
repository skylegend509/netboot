package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"netboot/internal/config"
	"netboot/internal/service"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func prepareTestEnv(t *testing.T) (*Handler, string) {
	tmpDir, err := os.MkdirTemp("", "netboot_handler_test")
	if err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		ISODir: tmpDir,
		Port:   ":8080",
	}
	isoSvc := service.NewISOService(cfg)
	ipxeSvc := service.NewIPXEService(cfg)
	h := NewHandler(cfg, isoSvc, ipxeSvc)

	return h, tmpDir
}

func TestHandleListISOs(t *testing.T) {
	h, tmpDir := prepareTestEnv(t)
	defer os.RemoveAll(tmpDir)

	req := httptest.NewRequest(http.MethodGet, "/api/isos", nil)
	w := httptest.NewRecorder()

	h.HandleListISOs(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %v", res.Status)
	}

	var isos []interface{}
	if err := json.NewDecoder(res.Body).Decode(&isos); err != nil {
		t.Fatal(err)
	}

	// Should be empty initially
	if len(isos) != 0 {
		t.Errorf("expected 0 isos, got %d", len(isos))
	}
}

func TestHandleUpload(t *testing.T) {
	h, tmpDir := prepareTestEnv(t)
	defer os.RemoveAll(tmpDir)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.iso")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("iso content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.HandleUpload(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %v", res.Status)
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(tmpDir, "test.iso")); os.IsNotExist(err) {
		t.Error("file was not uploaded")
	}
}

func TestHandleDelete(t *testing.T) {
	h, tmpDir := prepareTestEnv(t)
	defer os.RemoveAll(tmpDir)

	// Create dummy file
	os.WriteFile(filepath.Join(tmpDir, "delete_me.iso"), []byte("data"), 0644)
	// Force sync
	h.isoService.Sync()

	req := httptest.NewRequest(http.MethodDelete, "/api/delete?name=delete_me.iso", nil)
	w := httptest.NewRecorder()

	h.HandleDelete(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %v", res.Status)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "delete_me.iso")); !os.IsNotExist(err) {
		t.Error("file was not deleted")
	}
}

func TestHandleIPXE(t *testing.T) {
	h, tmpDir := prepareTestEnv(t)
	defer os.RemoveAll(tmpDir)

	// Create dummy ISO
	os.WriteFile(filepath.Join(tmpDir, "boot.iso"), []byte("boot"), 0644)
	h.isoService.Sync()

	req := httptest.NewRequest(http.MethodGet, "/boot.ipxe", nil)
	w := httptest.NewRecorder()

	h.HandleIPXE(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %v", res.Status)
	}

	body := w.Body.String()
	if !strings.Contains(body, "#!ipxe") {
		t.Error("response does not look like iPXE script")
	}
	if !strings.Contains(body, "boot.iso") {
		t.Error("response does not contain boot.iso")
	}
}
