package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Port           string
	ISODir         string
	UploadDir      string
	AdminToken     string
	MaxUploadBytes int64
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", ":8080"),
		ISODir:         getEnv("ISO_DIR", "./isos"),
		UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		AdminToken:     getEnv("ADMIN_TOKEN", ""),
		MaxUploadBytes: 10 << 30,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (c *Config) AbsISODir() string {
	abs, _ := filepath.Abs(c.ISODir)
	return abs
}
