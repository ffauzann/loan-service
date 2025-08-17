package vault

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/config/internal"
	vault "github.com/hashicorp/vault/api"
)

type Config struct {
	addr       string
	token      string
	mountPath  string
	secretPath string
	data       map[string]interface{}
	client     *vault.Client
	lock       *sync.RWMutex
	interval   time.Duration
	ticker     *time.Ticker
	close      chan struct{}
	watchers   []internal.Watcher
	kvversion  uint8
}

func (c *Config) reload() {
	resp, err := read(c.client, fmt.Sprintf("%s%s", c.mountPath, c.secretPath))
	if err != nil {
		log.Print(err)
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.kvversion == 2 {
		var ok bool
		c.data, ok = resp.Data["data"].(map[string]interface{})
		if !ok {
			c.data = resp.Data
		}
	} else {
		c.data = resp.Data
	}
	for _, watcher := range c.watchers {
		watcher.Update(c)
	}
}

func (c *Config) startTimer() {
	c.close = make(chan struct{})
	c.ticker = time.NewTicker(c.interval)
	go func() {
		for {
			select {
			case <-c.ticker.C:
				c.reload()
			case <-c.close:
				c.ticker.Stop()
				return
			}
		}
	}()
}

// SetFromMap sets secret from given data.
func SetFromMap(data map[string]interface{}) (*Config, error) {
	return &Config{
		data: data,
	}, nil
}

func (c *Config) Close() error {
	if c.interval > 0 {
		close(c.close)
	}

	for _, watcher := range c.watchers {
		watcher.Close()
	}

	return nil
}

func (c *Config) GetAll() map[string]interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data
}

func (c *Config) extractValue(key string) interface{} {
	keys := strings.Split(key, ".")
	if len(keys) <= 1 {
		return c.data[key]
	}

	parent, ok := c.data[keys[0]].(map[string]interface{})
	if !ok {
		return c.data[key]
	}

	for i := 1; i < len(keys)-1; i++ {
		parent, ok = parent[keys[i]].(map[string]interface{})
		if !ok {
			return parent[keys[len(keys)-1]]
		}
	}

	return parent[keys[len(keys)-1]]
}

func (c *Config) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.extractValue(key)
}

func (c *Config) GetInt(key string) int64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.extractValue(key).(int64)
	if ok {
		return value
	}
	return 0
}

func (c *Config) GetString(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.extractValue(key).(string)
	if ok {
		return value
	}
	return ""
}

func (c *Config) GetBool(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.extractValue(key).(bool)
	if ok {
		return value
	}
	return false
}

func (c *Config) GetFloat(key string) float64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.extractValue(key).(float64)
	if ok {
		return value
	}
	return 0
}

func (c *Config) GetBinary(key string) []byte {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.extractValue(key).(string)
	if ok {
		bytes, err := base64.StdEncoding.DecodeString(value)
		if err == nil {
			return bytes
		}
	}
	return nil
}

func (c *Config) GetArray(key string) []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.extractValue(key).([]string)
	if ok {
		return value
	}
	return nil
}

func (c *Config) GetMap(key string) map[string]string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.extractValue(key).(map[string]string)
	if ok {
		return value
	}
	return nil
}

func (c *Config) Watch(keys ...string) <-chan []string {
	watcher := internal.NewWatcher(keys, c)
	c.watchers = append(c.watchers, watcher)
	return watcher.Change()
}
