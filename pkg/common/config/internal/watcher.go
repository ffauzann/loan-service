package internal

import (
	"github.com/ffauzann/loan-service/pkg/common/config"
)

type Watcher struct {
	configMap map[string]interface{}
	channel   chan []string
}

func NewWatcher(keys []string, cfg config.Config) Watcher {
	cfgMap := make(map[string]interface{})
	for _, key := range keys {
		cfgMap[key] = cfg.Get(key)
	}

	return Watcher{configMap: cfgMap, channel: make(chan []string)}
}

func (w Watcher) Change() <-chan []string {
	return w.channel
}

func (w Watcher) Update(cfg config.Config) {
	change := false
	var keys []string
	for key, val := range w.configMap {
		if newVal := cfg.Get(key); newVal != val {
			w.configMap[key] = newVal
			keys = append(keys, key)
			change = true
		}
	}

	if change {
		w.channel <- keys
	}
}

func (w Watcher) Close() {
	close(w.channel)
}
