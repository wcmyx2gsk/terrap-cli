package providers

import (
	"fmt"
	"strings"
)

// Provider represents a Terraform provider with its name and version constraints.
type Provider struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Version string `json:"version"`
}

// String returns a human-readable representation of the provider.
func (p Provider) String() string {
	return fmt.Sprintf("%s (%s) @ %s", p.Name, p.Source, p.Version)
}

// IsValid checks that the provider has the minimum required fields.
// Note: Version is intentionally not required here since some providers
// may rely on implicit version constraints defined elsewhere.
func (p Provider) IsValid() bool {
	return strings.TrimSpace(p.Name) != "" && strings.TrimSpace(p.Source) != ""
}

// NormalizedSource returns the source in lowercase for consistent comparison.
func (p Provider) NormalizedSource() string {
	return strings.ToLower(strings.TrimSpace(p.Source))
}

// ProviderList is a slice of Provider with helper methods.
type ProviderList []Provider

// FindByName returns the first provider matching the given name, or false if not found.
func (pl ProviderList) FindByName(name string) (Provider, bool) {
	for _, p := range pl {
		if strings.EqualFold(p.Name, name) {
			return p, true
		}
	}
	return Provider{}, false
}

// Names returns a slice of all provider names.
func (pl ProviderList) Names() []string {
	names := make([]string, 0, len(pl))
	for _, p := range pl {
		names = append(names, p.Name)
	}
	return names
}

// FilterValid returns only providers that pass the IsValid check.
func (pl ProviderList) FilterValid() ProviderList {
	valid := make(ProviderList, 0, len(pl))
	for _, p := range pl {
		if p.IsValid() {
			valid = append(valid, p)
		}
	}
	return valid
}
