package config

import (
	"os"
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
	return os.ReadFile(name)
}

// NewEnvironment ...
func NewEnvironment() OS {
	return &customOS{}
}
