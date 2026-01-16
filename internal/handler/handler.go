package handler

import (
	"encoding/json"
	"net/http"
	"netboot/internal/config"
	"netboot/internal/service"
	"strings"
)

type Handler struct {
	cfg         *config.Config
	isoService  *service.ISOService
	ipxeService *service.IPXEService
}

func NewHandler(cfg *config.Config, isoSvc *service.ISOService, ipxeSvc *service.IPXEService) *Handler {
	return &Handler{
		cfg:         cfg,
		isoService:  isoSvc,
		ipxeService: ipxeSvc,
	}
}

func (h *Handler) HandleListISOs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	isos := h.isoService.List()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(isos)
}

func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit upload size to 10GB for example, or leave unlimited
	r.ParseMultipartForm(10 << 30) // 10 GB

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".iso") {
		http.Error(w, "Only ISO files allowed", http.StatusBadRequest)
		return
	}

	if err := h.isoService.Add(file, header); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	if err := h.isoService.Delete(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleIPXE(w http.ResponseWriter, r *http.Request) {
	isos := h.isoService.List()
	script, err := h.ipxeService.GenerateScript(isos)
	if err != nil {
		http.Error(w, "Failed to generate script", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(script))
}
