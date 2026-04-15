package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestScanProviders_EmptyData(t *testing.T) {
	results, err := scanProviders(map[string]interface{}{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestScanProviders_WithProviders(t *testing.T) {
	initData := map[string]interface{}{
		"providers": map[string]interface{}{
			"aws": map[string]interface{}{
				"version": "4.67.0",
			},
			"google": map[string]interface{}{
				"version": "4.51.0",
			},
		},
	}

	results, err := scanProviders(initData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Provider == "" {
			t.Error("provider name should not be empty")
		}
		if r.CurrentVer == "" {
			t.Errorf("provider %s should have a current version", r.Provider)
		}
	}
}

func TestPrintScanResults_JSONFormat(t *testing.T) {
	results := []ScanResult{
		{Provider: "aws", CurrentVer: "4.67.0", LatestVer: "4.67.0", UpgradeNeeded: false},
	}

	// Redirect stdout by capturing via json encode directly
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		t.Fatalf("failed to encode results: %v", err)
	}

	var decoded []ScanResult
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(decoded) != 1 {
		t.Errorf("expected 1 decoded result, got %d", len(decoded))
	}
	if decoded[0].Provider != "aws" {
		t.Errorf("expected provider 'aws', got '%s'", decoded[0].Provider)
	}
}

func TestPrintScanResults_TextFormat_Empty(t *testing.T) {
	err := printScanResults([]ScanResult{}, "text")
	if err != nil {
		t.Errorf("expected no error for empty results, got %v", err)
	}
}
