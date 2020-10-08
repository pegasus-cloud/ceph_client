package utility

import (
	"time"

	"github.com/bluele/gcache"
)

// GCacheBuilder is defined about a simple cache configuraion.
// It contains cache name and cache size.
type GCacheBuilder struct {
	name   string
	size   int
	expire time.Duration
}

var (
	caches map[string]gcache.Cache
)

func init() {
	caches = make(map[string]gcache.Cache)
}

// New return a cache builder accroding to configuration.
func New(cacheName string, size int) *GCacheBuilder {
	return &GCacheBuilder{
		name: cacheName,
		size: size,
	}
}

// NewWithExpire return a cache builder accroding to configuration.
func NewWithExpire(cacheName string, size int, expire time.Duration) *GCacheBuilder {
	return &GCacheBuilder{
		name:   cacheName,
		size:   size,
		expire: expire,
	}
}

// Build is used to build a cache entity.
// It's base on LRU policy and can be used anywhere in this project.
func (cb *GCacheBuilder) Build() {
	if cb.size <= 0 {
		return
	}
	caches[cb.name] = gcache.New(cb.size).
		LRU().
		Build()
}

//BuildWithExpire ...
func (cb *GCacheBuilder) BuildWithExpire() {
	if cb.size <= 0 {
		return
	}
	caches[cb.name] = gcache.New(cb.size).
		LRU().
		Expiration(cb.expire).
		Build()
}

// UseCache is used to specify a cache by cache name.
func UseCache(cacheName string) gcache.Cache {
	return caches[cacheName]
}
