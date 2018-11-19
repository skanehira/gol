package config

import "runtime"

type Config struct {
	OS              string
	ApplicationPath string
}

func New() *Config {
	var path string

	switch runtime.GOOS {
	case "darwin":
		path = "/Applications"
	case "linux":
		path = "/opt"
	}

	return &Config{
		ApplicationPath: path,
		OS:              runtime.GOOS,
	}
}

func (c *Config) IsDarwin() bool {
	return c.OS == "darwin"
}

func (c *Config) IsLinux() bool {
	return c.OS == "linux"
}
