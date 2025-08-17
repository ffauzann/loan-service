package vault

import (
	"os"
	"sync"
	"time"
)

type Option func(*Config)

func readOptions(path string, opts ...Option) *Config {
	// Default KV version is 1. It'll change if we specify the WithKVVersion() in option
	vaultConfig := &Config{
		addr:       os.Getenv("VAULT_ADDR"),
		token:      os.Getenv("VAULT_TOKEN"),
		mountPath:  os.Getenv("VAULT_MOUNT_PATH"),
		secretPath: path,
		lock:       &sync.RWMutex{},
		kvversion:  1,
	}

	// Apply all option to vaultConfig
	for _, opt := range opts {
		opt(vaultConfig)
	}

	return vaultConfig
}

// WithAddress sets specific vault address
func WithAddress(addr string) Option {
	return func(cfg *Config) {
		cfg.addr = addr
	}
}

// WithToken sets specific vault token
func WithToken(token string) Option {
	return func(cfg *Config) {
		cfg.token = token
	}
}

// WithMountPath sets specific vault token
func WithMountPath(mountPath string) Option {
	return func(cfg *Config) {
		cfg.mountPath = mountPath
	}
}

// WithReloadInterval sets interval duration to reload Config from vault
func WithReloadInterval(interval time.Duration) Option {
	return func(cfg *Config) {
		if interval > 0 && interval < time.Minute {
			interval = time.Minute
		}

		cfg.interval = interval

		if interval > 0 {
			cfg.startTimer()
		}
	}
}

// WithKVVersion specify the KV Version.
func WithKVVersion(version uint8) Option {
	return func(cfg *Config) {
		if version > 0 && version <= 2 {
			cfg.kvversion = version
		} else {
			cfg.kvversion = 1
		}
	}
}
