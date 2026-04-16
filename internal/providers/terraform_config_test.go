package providers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTerrapConfig(t *testing.T, dir string, cfg *TerraformConfigFile) {
	t.Helper()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".terrap.json"), data, 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func TestReadTerraformConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	cfg := &TerraformConfigFile{
		Terraform: TerraformConfig{
			RequiredProviders: map[string]ProviderRequirement{
				"aws": {Source: "hashicorp/aws", Version: "~> 4.0"},
			},
		},
	}
	writeTerrapConfig(t, dir, cfg)

	got, err := ReadTerraformConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req, ok := got.Terraform.RequiredProviders["aws"]; !ok {
		t.Fatal("expected aws provider in config")
	} else if req.Source != "hashicorp/aws" {
		t.Errorf("expected source hashicorp/aws, got %s", req.Source)
	} else if req.Version != "~> 4.0" {
		// Also verify the version constraint is preserved correctly
		t.Errorf("expected version ~> 4.0, got %s", req.Version)
	}
}

func TestReadTerraformConfig_Missing(t *testing.T) {
	dir := t.TempDir()
	_, err := ReadTerraformConfig(dir)
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestWriteTerraformConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfg := &TerraformConfigFile{
		Terraform: TerraformConfig{
			RequiredProviders: map[string]ProviderRequirement{
				"google":  {Source: "hashicorp/google", Version: ">= 3.0"},
				"azurerm": {Source: "hashicorp/azurerm", Version: "~> 2.99"},
			},
		},
	}
	if err := WriteTerraformConfig(dir, cfg); err != nil {
		t.Fatalf("write error: %v", err)
	}
	got, err := ReadTerraformConfig(dir)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if len(got.Terraform.RequiredProviders) != 2 {
		t.Errorf("expected 2 providers, got %d", len(got.Terraform.RequiredProviders))
	}
}

func TestTerraformConfigFile_ToProviderList(t *testing.T) {
	cfg := &TerraformConfigFile{
		Terraform: TerraformConfig{
			RequiredProviders: map[string]ProviderRequirement{
				"aws":    {Source: "hashicorp/aws", Version: "~> 4.0"},
				"random": {Source: "hashicorp/random", Version: ">= 3.1"},
			},
		},
	}
	list := cfg.ToProviderList()
	if len(list) != 2 {
		t.Errorf("expected 2 providers, got %d", len(list))
	}
	names := list.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}
