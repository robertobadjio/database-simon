package config

import (
	"os"
	"path/filepath"
)

// OS ...
type OS interface {
	GetEnv(string) string
	ReadFile(string) ([]byte, error)
}

type customOS struct{}

// GetEnv ...
func (customOS) GetEnv(key string) string {
	return os.Getenv(key)
}

// ReadFile ...
func (customOS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(name))
}

// NewEnvironment ...
func NewEnvironment() OS {
	return &customOS{}
}
