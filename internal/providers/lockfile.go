package providers

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const lockFileName = ".terraform.lock.hcl.json"

// LockEntry represents a single provider entry in the lock file.
type LockEntry struct {
	Source  string `json:"source"`
	Version string `json:"version"`
	Hash    string `json:"hash,omitempty"`
}

// LockFile represents the parsed terraform lock file state.
type LockFile struct {
	Providers map[string]LockEntry `json:"providers"`
}

// NewLockFile creates an empty LockFile.
func NewLockFile() *LockFile {
	return &LockFile{
		Providers: make(map[string]LockEntry),
	}
}

// ReadLockFile reads and parses a lock file from the given directory.
func ReadLockFile(dir string) (*LockFile, error) {
	path := filepath.Join(dir, lockFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewLockFile(), nil
		}
		return nil, err
	}
	lf := NewLockFile()
	if err := json.Unmarshal(data, lf); err != nil {
		return nil, err
	}
	return lf, nil
}

// WriteLockFile serializes and writes the lock file to the given directory.
func WriteLockFile(dir string, lf *LockFile) error {
	path := filepath.Join(dir, lockFileName)
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Upsert adds or updates a provider entry in the lock file.
func (lf *LockFile) Upsert(source, version, hash string) {
	lf.Providers[source] = LockEntry{
		Source:  source,
		Version: version,
		Hash:    hash,
	}
}

// Remove deletes a provider entry by source.
func (lf *LockFile) Remove(source string) {
	delete(lf.Providers, source)
}

// ToProviderList converts lock file entries into a ProviderList.
func (lf *LockFile) ToProviderList() ProviderList {
	list := make(ProviderList, 0, len(lf.Providers))
	for _, entry := range lf.Providers {
		list = append(list, Provider{
			Source:  entry.Source,
			Version: entry.Version,
		})
	}
	return list
}
