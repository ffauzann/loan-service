package config

import "io"

// Config is interface to get configuration value
type Config interface {
	io.Closer

	// GetAll returns all configuration value
	GetAll() map[string]interface{}

	// Get returns configuration value as interface
	Get(string) interface{}

	// GetInt returns configuration value as integer 64 bit
	GetInt(string) int64

	// GetString returns configuration value as string
	GetString(string) string

	// GetBool returns configuration value as boolean
	GetBool(string) bool

	// GetFloat returns configuration value as float 64 bit
	GetFloat(string) float64

	// GetBinary returns configuration value as byte array,
	// configuration value is stored as base64 encoded
	GetBinary(string) []byte

	// GetArray returns configuration value as array
	// configuration value is stored with format <element1>,<element2>,...
	GetArray(string) []string

	// GetMap returns configuration value as map
	// configuration value is stored with format <key1>:<value1>,<key2>:<value2>,...
	GetMap(string) map[string]string

	// Watch watches for value changes on specific keys
	// The channel will return keys when the values are changed
	Watch(...string) <-chan []string
}
