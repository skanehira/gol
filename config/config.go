package config

import "runtime"

const (
	MacOS   = "darwin"
	Linux   = "linux"
	Windows = "windows"
)

type Config struct {
	OS              string
	ApplicationPath []string
}

func New() *Config {
	var conf = &Config{
		OS: runtime.GOOS,
	}

	switch runtime.GOOS {
	case MacOS:
		conf.ApplicationPath = append(conf.ApplicationPath, []string{"/Applications", "/System/Library/CoreServices/Applications"}...)
	case Linux:
		// TODO support linux
	case Windows:
		// TODO support windows
	}

	return conf
}
