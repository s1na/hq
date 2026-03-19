package cache

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Cache struct {
	Dir     string
	Enabled bool
}

func New(dir string, enabled bool) (*Cache, error) {
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(home, ".cache", "hq")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Cache{Dir: dir, Enabled: enabled}, nil
}

func (c *Cache) key(url string) string {
	h := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x", h)
}

// Get returns cached data if it exists and hasn't expired.
// ttl <= 0 means the entry never expires (immutable content).
func (c *Cache) Get(url string, ttl time.Duration) ([]byte, bool) {
	if !c.Enabled {
		return nil, false
	}
	path := filepath.Join(c.Dir, c.key(url))
	info, err := os.Stat(path)
	if err != nil {
		return nil, false
	}
	if ttl > 0 && time.Since(info.ModTime()) > ttl {
		return nil, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return data, true
}

// Put stores data in the cache. Always writes, even if cache reads are disabled.
func (c *Cache) Put(url string, data []byte) error {
	path := filepath.Join(c.Dir, c.key(url))
	return os.WriteFile(path, data, 0644)
}
