package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const registryBaseURL = "https://registry.terraform.io/v1/providers"

// RegistryResponse represents the response from the Terraform registry API.
type RegistryResponse struct {
	Versions []RegistryVersion `json:"versions"`
}

// RegistryVersion holds a single version entry from the registry.
type RegistryVersion struct {
	Version string `json:"version"`
}

// RegistryClient fetches provider information from the Terraform registry.
type RegistryClient struct {
	httpClient *http.Client
}

// NewRegistryClient creates a new RegistryClient with a default timeout.
func NewRegistryClient() *RegistryClient {
	return &RegistryClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetLatestVersion queries the registry for the latest version of a provider.
func (rc *RegistryClient) GetLatestVersion(source string) (string, error) {
	url := fmt.Sprintf("%s/%s/versions", registryBaseURL, source)
	resp, err := rc.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("registry request failed for %s: %w", source, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("registry returned status %d for %s", resp.StatusCode, source)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read registry response: %w", err)
	}

	var result RegistryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse registry response: %w", err)
	}

	if len(result.Versions) == 0 {
		return "", fmt.Errorf("no versions found for provider %s", source)
	}

	return result.Versions[len(result.Versions)-1].Version, nil
}
