package model

// Reusable config goes here.
type AppConfig struct {
	Encryption Encryption
	Jwt        JwtConfig
	Auth       AuthConfig
	Dependency DependencyConfig
}

type Encryption struct {
	Cost uint8
}

type AuthConfig struct {
	ExcludedMethods []string
}

type DependencyConfig struct{}
