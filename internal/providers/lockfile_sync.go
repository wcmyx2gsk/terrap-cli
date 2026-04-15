package providers

import "fmt"

// SyncResult summarises the outcome of a lock file sync operation.
type SyncResult struct {
	Added   []string
	Updated []string
	Removed []string
}

// String returns a human-readable summary of the sync result.
func (r SyncResult) String() string {
	return fmt.Sprintf("added=%d updated=%d removed=%d",
		len(r.Added), len(r.Updated), len(r.Removed))
}

// HasChanges returns true if the sync result contains any modifications.
func (r SyncResult) HasChanges() bool {
	return len(r.Added)+len(r.Updated)+len(r.Removed) > 0
}

// SyncLockFile reconciles a LockFile against a desired ProviderList.
// Providers in desired but not in lf are added; providers present in lf
// but absent from desired are removed; version mismatches are updated.
func SyncLockFile(lf *LockFile, desired ProviderList) SyncResult {
	result := SyncResult{}

	desiredMap := make(map[string]Provider, len(desired))
	for _, p := range desired {
		desiredMap[p.NormalizedSource()] = p
	}

	// Detect removals
	for source := range lf.Providers {
		if _, ok := desiredMap[source]; !ok {
			lf.Remove(source)
			result.Removed = append(result.Removed, source)
		}
	}

	// Detect additions and updates
	for source, p := range desiredMap {
		existing, exists := lf.Providers[source]
		if !exists {
			lf.Upsert(source, p.Version, "")
			result.Added = append(result.Added, source)
		} else if existing.Version != p.Version {
			lf.Upsert(source, p.Version, "")
			result.Updated = append(result.Updated, source)
		}
	}

	return result
}
