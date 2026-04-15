package providers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLockFile_UpsertAndRemove(t *testing.T) {
	lf := NewLockFile()
	lf.Upsert("hashicorp/aws", "4.0.0", "h1:abc")
	lf.Upsert("hashicorp/google", "3.1.0", "")

	if len(lf.Providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(lf.Providers))
	}

	lf.Remove("hashicorp/aws")
	if _, ok := lf.Providers["hashicorp/aws"]; ok {
		t.Fatal("expected hashicorp/aws to be removed")
	}
}

func TestLockFile_WriteAndRead(t *testing.T) {
	dir := t.TempDir()
	lf := NewLockFile()
	lf.Upsert("hashicorp/aws", "4.0.0", "h1:abc")

	if err := WriteLockFile(dir, lf); err != nil {
		t.Fatalf("WriteLockFile error: %v", err)
	}

	path := filepath.Join(dir, lockFileName)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("lock file not written: %v", err)
	}

	loaded, err := ReadLockFile(dir)
	if err != nil {
		t.Fatalf("ReadLockFile error: %v", err)
	}
	entry, ok := loaded.Providers["hashicorp/aws"]
	if !ok {
		t.Fatal("expected hashicorp/aws in loaded lock file")
	}
	if entry.Version != "4.0.0" {
		t.Errorf("expected version 4.0.0, got %s", entry.Version)
	}
	if entry.Hash != "h1:abc" {
		t.Errorf("expected hash h1:abc, got %s", entry.Hash)
	}
}

func TestReadLockFile_MissingFile(t *testing.T) {
	dir := t.TempDir()
	lf, err := ReadLockFile(dir)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(lf.Providers) != 0 {
		t.Errorf("expected empty providers map, got %d entries", len(lf.Providers))
	}
}

func TestLockFile_ToProviderList(t *testing.T) {
	lf := NewLockFile()
	lf.Upsert("hashicorp/aws", "4.0.0", "")
	lf.Upsert("hashicorp/random", "3.0.0", "")

	list := lf.ToProviderList()
	if len(list) != 2 {
		t.Errorf("expected 2 providers in list, got %d", len(list))
	}
}
