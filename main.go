package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// side-effect imports register sub-commands via cobra init() hooks
	_ "github.com/terrap-cli/terrap-cli/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "terrap",
	Short: "terrap — Terraform provider version manager",
	Long: `terrap helps you keep Terraform provider versions up to date.

It can scan your configuration, check for newer provider versions,
generate upgrade plans, and synchronise your lock file.

Run 'terrap <command> --help' for details on a specific command.`,
	// Silence default error printing so we can format it ourselves.
	SilenceErrors: true,
	SilenceUsage:  true,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
