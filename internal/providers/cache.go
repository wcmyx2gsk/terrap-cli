package providers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 24 * time.Hour

// CacheEntry holds a cached version check result with a timestamp.
type CacheEntry struct {
	Result    VersionCheckResult `json:"result"`
	CachedAt  time.Time          `json:"cached_at"`
}

// ProviderCache manages on-disk caching of version check results.
type ProviderCache struct {
	cacheDir string
}

// NewProviderCache creates a new ProviderCache rooted at cacheDir.
func NewProviderCache(cacheDir string) *ProviderCache {
	return &ProviderCache{cacheDir: cacheDir}
}

// cacheFilePath returns the path for a provider's cache file.
func (c *ProviderCache) cacheFilePath(providerName string) string {
	return filepath.Join(c.cacheDir, providerName+".json")
}

// Get retrieves a cached VersionCheckResult if it exists and is still valid.
func (c *ProviderCache) Get(providerName string) (*VersionCheckResult, bool) {
	data, err := os.ReadFile(c.cacheFilePath(providerName))
	if err != nil {
		return nil, false
	}
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}
	if time.Since(entry.CachedAt) > cacheTTL {
		return nil, false
	}
	return &entry.Result, true
}

// Set writes a VersionCheckResult to the cache.
func (c *ProviderCache) Set(providerName string, result VersionCheckResult) error {
	if err := os.MkdirAll(c.cacheDir, 0755); err != nil {
		return err
	}
	entry := CacheEntry{
		Result:   result,
		CachedAt: time.Now(),
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(c.cacheFilePath(providerName), data, 0644)
}

// Invalidate removes the cache entry for the given provider.
func (c *ProviderCache) Invalidate(providerName string) error {
	err := os.Remove(c.cacheFilePath(providerName))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
