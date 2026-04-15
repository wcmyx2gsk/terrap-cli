package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TerraformConfig represents a parsed terraform block from a configuration file.
type TerraformConfig struct {
	RequiredProviders map[string]ProviderRequirement `json:"required_providers"`
}

// ProviderRequirement holds the source and version constraint for a required provider.
type ProviderRequirement struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

// TerraformConfigFile represents the top-level structure of a .terraform.lock.hcl
// adjacent JSON metadata file used by terrap-cli.
type TerraformConfigFile struct {
	Terraform TerraformConfig `json:"terraform"`
}

// ReadTerraformConfig reads and parses a terrap JSON config file from the given directory.
func ReadTerraformConfig(dir string) (*TerraformConfigFile, error) {
	path := filepath.Join(dir, ".terrap.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("terrap config not found at %s: run 'terrap init' first", path)
		}
		return nil, fmt.Errorf("reading terrap config: %w", err)
	}

	var cfg TerraformConfigFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing terrap config: %w", err)
	}
	return &cfg, nil
}

// WriteTerraformConfig writes the given config to .terrap.json in dir.
func WriteTerraformConfig(dir string, cfg *TerraformConfigFile) error {
	path := filepath.Join(dir, ".terrap.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling terrap config: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing terrap config: %w", err)
	}
	return nil
}

// ToProviderList converts the required_providers map into a ProviderList.
func (c *TerraformConfigFile) ToProviderList() ProviderList {
	var list ProviderList
	for name, req := range c.Terraform.RequiredProviders {
		list = append(list, Provider{
			Name:    name,
			Source:  req.Source,
			Version: req.Version,
		})
	}
	return list
}
