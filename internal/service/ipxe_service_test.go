package service

import (
	"netboot/internal/config"
	"netboot/internal/model"
	"strings"
	"testing"
)

func TestIPXEService_GenerateScript(t *testing.T) {
	cfg := &config.Config{
		Port: ":9090",
	}
	svc := NewIPXEService(cfg)

	isos := []model.ISO{
		{Name: "ubuntu.iso", Path: "ubuntu.iso", Size: 1024},
		{Name: "arch.iso", Path: "arch.iso", Size: 2048},
	}

	script, err := svc.GenerateScript(isos)
	if err != nil {
		t.Fatalf("GenerateScript failed: %v", err)
	}

	// Verify crucial parts of the script
	expectedStrings := []string{
		"#!ipxe",
		"item iso0 ubuntu.iso",
		"item iso1 arch.iso",
		"kernel http://${next-server}:9090/isos/ubuntu.iso",
		"kernel http://${next-server}:9090/isos/arch.iso",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(script, s) {
			t.Errorf("Script missing expected content: %s", s)
		}
	}
}
