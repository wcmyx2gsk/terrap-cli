package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		p        Provider
		expected bool
	}{
		{"valid provider", Provider{Name: "aws", Source: "hashicorp/aws", Version: "~> 4.0"}, true},
		{"missing name", Provider{Name: "", Source: "hashicorp/aws"}, false},
		{"missing source", Provider{Name: "aws", Source: ""}, false},
		{"whitespace name", Provider{Name: "  ", Source: "hashicorp/aws"}, false},
		// also check whitespace-only source, similar to whitespace name case
		{"whitespace source", Provider{Name: "aws", Source: "  "}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.p.IsValid())
		})
	}
}

func TestProvider_NormalizedSource(t *testing.T) {
	p := Provider{Name: "AWS", Source: "  HashiCorp/AWS  "}
	assert.Equal(t, "hashicorp/aws", p.NormalizedSource())
}

func TestProviderList_FindByName(t *testing.T) {
	list := ProviderList{
		{Name: "aws", Source: "hashicorp/aws"},
		{Name: "google", Source: "hashicorp/google"},
	}

	p, ok := list.FindByName("aws")
	assert.True(t, ok)
	assert.Equal(t, "aws", p.Name)

	_, ok = list.FindByName("azure")
	assert.False(t, ok)
}

func TestProviderList_Names(t *testing.T) {
	list := ProviderList{
		{Name: "aws", Source: "hashicorp/aws"},
		{Name: "google", Source: "hashicorp/google"},
	}
	names := list.Names()
	assert.ElementsMatch(t, []string{"aws", "google"}, names)
}
