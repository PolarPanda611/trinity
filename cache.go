package trinity

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
)

func initDefaultCache(cacheAlgorithm string, cacheSize int, timeout int) gcache.Cache {
	fmt.Println(cacheAlgorithm)
	switch cacheAlgorithm {
	case "LFU":
		return gcache.New(cacheSize).LFU().Expiration(time.Duration(timeout) * time.Hour).Build()
	case "LRU":
		return gcache.New(cacheSize).LRU().Expiration(time.Duration(timeout) * time.Hour).Build()
	case "ARC":
		return gcache.New(cacheSize).ARC().Expiration(time.Duration(timeout) * time.Hour).Build()
	default:
		return gcache.New(cacheSize).Expiration(time.Duration(timeout) * time.Hour).Build()
	}

}

func (t *Trinity) initCache() {
	t.cache = initDefaultCache("LRU", t.setting.Cache.Gcache.CacheSize, t.setting.Cache.Gcache.Timeout)
}

// CleanCache clean the current cache
func (t *Trinity) CleanCache() {
	t.mu.Lock()
	t.cache = initDefaultCache("LRU", t.setting.Cache.Gcache.CacheSize, t.setting.Cache.Gcache.Timeout)
	t.mu.Unlock()
}

// GetCache  get vcfg
func (t *Trinity) GetCache() gcache.Cache {
	t.mu.RLock()
	c := t.cache
	t.mu.RUnlock()
	return c
}

// SetCache  get vcfg
func (t *Trinity) SetCache(cache gcache.Cache) *Trinity {
	t.mu.Lock()
	t.cache = cache
	t.mu.Unlock()
	return t
}
