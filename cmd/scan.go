package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type ScanResult struct {
	Provider    string            `json:"provider"`
	CurrentVer  string            `json:"current_version"`
	LatestVer   string            `json:"latest_version"`
	UpgradeNeeded bool            `json:"upgrade_needed"`
	Changes     []string          `json:"breaking_changes,omitempty"`
}

var outputFormat string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan terraform configurations for provider version drift",
	Long:  `Scan reads the local .terrap state and compares provider versions against the latest available releases.`,
	RunE:  runScan,
}

func runScan(cmd *cobra.Command, args []string) error {
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	statePath := filepath.Join(workDir, ".terrap", "init_data.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		return fmt.Errorf("no terrap state found. Run 'terrap init' first: %w", err)
	}

	var initData map[string]interface{}
	if err := json.Unmarshal(data, &initData); err != nil {
		return fmt.Errorf("failed to parse terrap state: %w", err)
	}

	results, err := scanProviders(initData)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	return printScanResults(results, outputFormat)
}

func scanProviders(initData map[string]interface{}) ([]ScanResult, error) {
	var results []ScanResult
	providers, ok := initData["providers"].(map[string]interface{})
	if !ok {
		return results, nil
	}

	for name, meta := range providers {
		providerMeta, ok := meta.(map[string]interface{})
		if !ok {
			continue
		}
		current, _ := providerMeta["version"].(string)
		results = append(results, ScanResult{
			Provider:      name,
			CurrentVer:    current,
			LatestVer:     current, // placeholder: real impl would fetch from registry
			UpgradeNeeded: false,
		})
	}
	return results, nil
}

func printScanResults(results []ScanResult, format string) error {
	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}
	if len(results) == 0 {
		fmt.Println("No providers found in terrap state.")
		return nil
	}
	for _, r := range results {
		status := "✔ up-to-date"
		if r.UpgradeNeeded {
			status = "⚠ upgrade available"
		}
		fmt.Printf("[%s] %s  current: %s  latest: %s  %s\n", r.Provider, status, r.CurrentVer, r.LatestVer, "")
	}
	return nil
}

func init() {
	scanCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(scanCmd)
}
