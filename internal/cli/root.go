// Package cli provides the command-line interface for sopsctl.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/enbiyagoral/sopsctl/internal/config"
	"github.com/enbiyagoral/sopsctl/internal/selector"
	"github.com/enbiyagoral/sopsctl/internal/sops"
)

var (
	// Global flags
	cfgFile     string
	profileName string
	noFzf       bool
	dryRun      bool

	// Shared instances
	cfg      *config.Config
	builder  *sops.ArgsBuilder
	executor *sops.Executor
)

var rootCmd = &cobra.Command{
	Use:   "sopsctl",
	Short: "SOPS profile manager",
	Long: `sopsctl is a profile manager for SOPS.

Select a profile, and sopsctl builds the right SOPS arguments for you.

Profile selection priority:
  1. -p/--profile flag (explicit)
  2. Directory mapping (directories in config)
  3. default_profile (if set in config)
  4. Interactive fzf selection`,
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

		// Initialize shared instances
		builder = sops.NewArgsBuilder()
		executor = sops.NewExecutor(cfg.Settings.SOPSPath)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: ~/.config/sopsctl/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&profileName, "profile", "p", "", "profile to use (skip fzf selection)")
	rootCmd.PersistentFlags().BoolVar(&noFzf, "no-fzf", false, "disable fzf selection")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show command without executing")

	// Add subcommands
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(configCmd)
}

// Execute runs the CLI.
func Execute() error {
	return rootCmd.Execute()
}

// selectProfile handles profile selection logic.
// Priority: -p flag > directory mapping > default_profile > fzf
func selectProfile() (*config.Profile, error) {
	// 1. Check if profile was specified via flag
	if profileName != "" {
		return cfg.GetProfile(profileName)
	}

	// 2. Check for directory-based profile
	cwd, err := os.Getwd()
	if err == nil {
		if name, auto, found := cfg.ResolveProfile(cwd); found {
			profile, err := cfg.GetProfile(name)
			if err != nil {
				return nil, err
			}

			if auto {
				fmt.Fprintf(os.Stderr, "Using profile: %s (directory auto-select)\n", name)
				return profile, nil
			}
		}
	}

	// 3. Check for default profile
	if cfg.DefaultProfile != "" {
		profile, err := cfg.GetProfile(cfg.DefaultProfile)
		if err == nil {
			fmt.Fprintf(os.Stderr, "Using profile: %s (default)\n", cfg.DefaultProfile)
			return profile, nil
		}
	}

	// 4. Use fzf selection (if not disabled)
	if noFzf {
		return nil, fmt.Errorf("no profile specified (use -p or set default_profile in config)")
	}

	profiles := cfg.ListProfiles()
	if len(profiles) == 0 {
		return nil, fmt.Errorf("no profiles configured\nRun 'sopsctl profile add' to create one")
	}

	sel := selector.NewFZFSelector()
	return sel.Select(profiles)
}
