package service

import (
	"io"
	"mime/multipart"
	"netboot/internal/config"
	"netboot/internal/model"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ISOService struct {
	cfg  *config.Config
	mu   sync.RWMutex
	isos []model.ISO
}

func NewISOService(cfg *config.Config) *ISOService {
	s := &ISOService{
		cfg: cfg,
	}
	// Initial scan
	os.MkdirAll(cfg.ISODir, 0755)
	s.Sync()
	return s
}

func (s *ISOService) List() []model.ISO {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Return a copy to avoid race conditions if the caller modifies the slice
	result := make([]model.ISO, len(s.isos))
	copy(result, s.isos)
	return result
}

func (s *ISOService) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var newISOs []model.ISO
	err := filepath.Walk(s.cfg.ISODir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors accessing files
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".iso") {
			relPath, _ := filepath.Rel(s.cfg.ISODir, path)
			// Ensure forward slashes for URLs
			relPath = strings.ReplaceAll(relPath, "\\", "/")
			
			newISOs = append(newISOs, model.ISO{
				Name: info.Name(),
				Size: info.Size(),
				Path: relPath,
			})
		}
		return nil
	})

	s.isos = newISOs
	return err
}

func (s *ISOService) Add(file multipart.File, header *multipart.FileHeader) error {
	dstPath := filepath.Join(s.cfg.ISODir, header.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	return s.Sync()
}

func (s *ISOService) Delete(name string) error {
	// Basic validation to prevent traversal, though filepath.Join handles simple cases.
	// We only allow deleting files in the root of ISO dir for now for safety.
	baseName := filepath.Base(name)
	targetPath := filepath.Join(s.cfg.ISODir, baseName)
	
	if err := os.Remove(targetPath); err != nil {
		return err
	}
	
	return s.Sync()
}
