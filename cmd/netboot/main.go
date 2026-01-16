package main

import (
	"fmt"
	"log"
	"net/http"
	"netboot/internal/config"
	"netboot/internal/handler"
	"netboot/internal/service"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Initialize Services
	isoService := service.NewISOService(cfg)
	ipxeService := service.NewIPXEService(cfg)

	// 3. Initialize Handlers
	h := handler.NewHandler(cfg, isoService, ipxeService)

	// 4. Setup Router
	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("/api/isos", h.HandleListISOs)
	mux.HandleFunc("/api/upload", h.HandleUpload)
	mux.HandleFunc("/api/delete", h.HandleDelete)
	mux.HandleFunc("/boot.ipxe", h.HandleIPXE)

	// Static Files (Frontend)
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	// ISO File Serving
	isoFs := http.StripPrefix("/isos/", http.FileServer(http.Dir(cfg.ISODir)))
	mux.Handle("/isos/", isoFs)

	// 5. Start Server
	serverAddr := cfg.Port
	fmt.Printf("🚀 PXE Boot Server running on http://localhost%s\n", serverAddr)
	fmt.Printf("📁 Serving ISOs from: %s\n", cfg.AbsISODir())
	fmt.Printf("📡 iPXE Script URL: http://localhost%s/boot.ipxe\n", serverAddr)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatal(err)
	}
}
