package config

import "runtime"

const (
	MacOS   = "darwin"
	Linux   = "linux"
	Windows = "windows"
)

type Config struct {
	OS              string
	ApplicationPath string
}

func New() *Config {
	var path string

	switch runtime.GOOS {
	case MacOS:
		path = "/Applications"
	case Linux:
		path = "/opt"
	case Windows:
		path = "C:¥¥Program Files"
	}

	return &Config{
		ApplicationPath: path,
		OS:              runtime.GOOS,
	}
}
