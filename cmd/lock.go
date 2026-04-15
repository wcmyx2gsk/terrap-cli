package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirrend/terrap-cli/internal/providers"
	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Sync the provider lock file with the current terraform configuration",
	RunE:  runLock,
}

func runLock(cmd *cobra.Command, _ []string) error {
	dir, _ := cmd.Flags().GetString("dir")
	if dir == "" {
		dir = "."
	}
	jsonOut, _ := cmd.Flags().GetBool("json")

	provList, err := providers.LoadProvidersFromDir(dir)
	if err != nil {
		return fmt.Errorf("failed to load providers: %w", err)
	}

	lf, err := providers.ReadLockFile(dir)
	if err != nil {
		return fmt.Errorf("failed to read lock file: %w", err)
	}

	result := providers.SyncLockFile(lf, provList)

	if result.HasChanges() {
		if err := providers.WriteLockFile(dir, lf); err != nil {
			return fmt.Errorf("failed to write lock file: %w", err)
		}
	}

	if jsonOut {
		return printLockJSON(result)
	}
	printLockText(result)
	return nil
}

func printLockJSON(result providers.SyncResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

func printLockText(result providers.SyncResult) {
	if !result.HasChanges() {
		fmt.Println("Lock file is already up to date.")
		return
	}
	for _, s := range result.Added {
		fmt.Printf("  + added:   %s\n", s)
	}
	for _, s := range result.Updated {
		fmt.Printf("  ~ updated: %s\n", s)
	}
	for _, s := range result.Removed {
		fmt.Printf("  - removed: %s\n", s)
	}
	fmt.Printf("\nSync complete: %s\n", result.String())
}

func init() {
	lockCmd.Flags().String("dir", ".", "Directory containing terraform configuration")
	lockCmd.Flags().Bool("json", false, "Output results as JSON")
	rootCmd.AddCommand(lockCmd)
}
