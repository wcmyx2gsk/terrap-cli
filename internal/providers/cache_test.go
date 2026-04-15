package providers

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func tempCacheDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "terrap-cache-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestProviderCache_SetAndGet(t *testing.T) {
	cache := NewProviderCache(tempCacheDir(t))
	result := VersionCheckResult{
		ProviderName:   "hashicorp/aws",
		CurrentVersion: "4.0.0",
		LatestVersion:  "5.0.0",
		Outdated:       true,
	}
	if err := cache.Set("aws", result); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	got, ok := cache.Get("aws")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if got.ProviderName != result.ProviderName {
		t.Errorf("expected %s, got %s", result.ProviderName, got.ProviderName)
	}
	if got.LatestVersion != result.LatestVersion {
		t.Errorf("expected %s, got %s", result.LatestVersion, got.LatestVersion)
	}
}

func TestProviderCache_Miss(t *testing.T) {
	cache := NewProviderCache(tempCacheDir(t))
	_, ok := cache.Get("nonexistent")
	if ok {
		t.Error("expected cache miss for nonexistent provider")
	}
}

func TestProviderCache_Invalidate(t *testing.T) {
	cache := NewProviderCache(tempCacheDir(t))
	result := VersionCheckResult{ProviderName: "hashicorp/google"}
	_ = cache.Set("google", result)
	if err := cache.Invalidate("google"); err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}
	_, ok := cache.Get("google")
	if ok {
		t.Error("expected cache miss after invalidation")
	}
}

func TestProviderCache_Expiry(t *testing.T) {
	dir := tempCacheDir(t)
	cache := NewProviderCache(dir)
	result := VersionCheckResult{ProviderName: "hashicorp/azurerm"}
	_ = cache.Set("azurerm", result)

	entry := CacheEntry{
		Result:   result,
		CachedAt: time.Now().Add(-25 * time.Hour),
	}
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(cache.cacheFilePath("azurerm"), data, 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, ok := cache.Get("azurerm")
	if ok {
		t.Error("expected cache miss for expired entry")
	}
}

func TestProviderCache_InvalidateNonExistent(t *testing.T) {
	cache := NewProviderCache(tempCacheDir(t))
	if err := cache.Invalidate("does-not-exist"); err != nil {
		t.Errorf("expected no error invalidating nonexistent entry, got: %v", err)
	}
}
