package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"iptc/pkg/logging"
	"k8s.io/utils/lru"
)

type Cache struct {
	dir        string
	lru        *lru.Cache
	expiration time.Duration
	mu         sync.Mutex
}

// NewCache initializes a new Cache with the specified directory, size, and expiration duration.
func NewCache(dir string, size int, expiration time.Duration) (*Cache, error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	lruCache := lru.New(size)

	return &Cache{
		dir:        dir,
		lru:        lruCache,
		expiration: expiration,
	}, nil
}

// Get retrieves a cached item by key. Returns the cached data or an error if not found.
func (c *Cache) Get(key string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if path, found := c.lru.Get(key); found {
		return os.ReadFile(path.(string))
	}
	return nil, os.ErrNotExist
}

// Set adds a new item to the cache with the specified key and data.
func (c *Cache) Set(key string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	path := filepath.Join(c.dir, key)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create cache subdirectory: %w", err)
	}

	if err := os.WriteFile(path, data, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	c.lru.Add(key, path)
	go c.expireKey(key, path, c.expiration)
	return nil
}

// expireKey removes a cached item after the expiration duration.
func (c *Cache) expireKey(key, path string, duration time.Duration) {
	time.Sleep(duration)
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lru.Remove(key)
	if err := os.Remove(path); err != nil {
		logging.Error(fmt.Sprintf("failed to remove expired cache file: %v", err))
	}
}
