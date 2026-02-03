// Package cli provides the command-line interface for sopsctl.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/enbiyagoral/sopsctl/internal/config"
)

var (
	// Global flags
	cfgFile     string
	profileName string

	// Shared instances
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "sopsctl",
	Short: "SOPS profile manager",
	Long: `sopsctl is a profile manager for SOPS.

Manage multiple age key files as profiles and set SOPS_AGE_KEY_FILE
for your shell, so you can use the native sops CLI directly.

Quick start:
  1. Initialize:     sopsctl config init
  2. Add profile:    sopsctl profile add stg --age-key-file ~/.sops/stg.txt
  3. Shell setup:    echo 'eval "$(sopsctl init zsh)"' >> ~/.zshrc
  4. Use profile:    sopsctl profile use stg
  5. Use sops:       sops -e -i secrets.yaml`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for certain commands
		if cmd.Name() == "init" || cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}

		// Load configuration
		path := cfgFile
		if path == "" {
			var err error
			path, err = config.DefaultConfigPath()
			if err != nil {
				return err
			}
		}

		var err error
		cfg, err = config.Load(path)
		if err != nil {
			// Allow profile and config commands without existing config
			if cmd.Parent() != nil && (cmd.Parent().Name() == "config" || cmd.Parent().Name() == "profile") {
				cfg = config.NewConfig()
				return nil
			}
			return fmt.Errorf("failed to load config: %w\nRun 'sopsctl config init' to create one", err)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: ~/.config/sopsctl/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&profileName, "profile", "p", "", "profile to use")

	// Add subcommands
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(configCmd)
}

// Execute runs the CLI.
func Execute() error {
	return rootCmd.Execute()
}
