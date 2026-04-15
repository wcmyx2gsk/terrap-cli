package providers

import "testing"

func TestSyncLockFile_AddNewProvider(t *testing.T) {
	lf := NewLockFile()
	desired := ProviderList{
		{Source: "hashicorp/aws", Version: "4.0.0"},
	}
	res := SyncLockFile(lf, desired)
	if len(res.Added) != 1 || res.Added[0] != "registry.terraform.io/hashicorp/aws" {
		t.Errorf("expected 1 addition, got %+v", res)
	}
	if res.HasChanges() == false {
		t.Error("expected HasChanges to be true")
	}
}

func TestSyncLockFile_RemoveProvider(t *testing.T) {
	lf := NewLockFile()
	lf.Upsert("registry.terraform.io/hashicorp/aws", "4.0.0", "")
	desired := ProviderList{}
	res := SyncLockFile(lf, desired)
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 removal, got %+v", res)
	}
	if _, ok := lf.Providers["registry.terraform.io/hashicorp/aws"]; ok {
		t.Error("provider should have been removed from lock file")
	}
}

func TestSyncLockFile_UpdateProvider(t *testing.T) {
	lf := NewLockFile()
	lf.Upsert("registry.terraform.io/hashicorp/aws", "3.0.0", "h1:old")
	desired := ProviderList{
		{Source: "hashicorp/aws", Version: "4.0.0"},
	}
	res := SyncLockFile(lf, desired)
	if len(res.Updated) != 1 {
		t.Errorf("expected 1 update, got %+v", res)
	}
	entry := lf.Providers["registry.terraform.io/hashicorp/aws"]
	if entry.Version != "4.0.0" {
		t.Errorf("expected updated version 4.0.0, got %s", entry.Version)
	}
}

func TestSyncLockFile_NoChanges(t *testing.T) {
	lf := NewLockFile()
	lf.Upsert("registry.terraform.io/hashicorp/aws", "4.0.0", "")
	desired := ProviderList{
		{Source: "hashicorp/aws", Version: "4.0.0"},
	}
	res := SyncLockFile(lf, desired)
	if res.HasChanges() {
		t.Errorf("expected no changes, got %+v", res)
	}
}

func TestSyncResult_String(t *testing.T) {
	r := SyncResult{Added: []string{"a"}, Updated: []string{"b", "c"}, Removed: []string{}}
	s := r.String()
	if s != "added=1 updated=2 removed=0" {
		t.Errorf("unexpected string: %s", s)
	}
}
