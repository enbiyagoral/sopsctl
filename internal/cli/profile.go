package cli

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/enbiyagoral/sopsy/internal/config"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage encryption profiles",
	Long:  `Manage SOPS encryption profiles (add, list, show, edit, remove).`,
}

var profileAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new profile",
	Long: `Add a new encryption profile.

Examples:
  # Add profile with age key file (recommended)
  sopsy profile add dev --description "Development" --age-key-file "~/.config/sops/age/keys.txt"
  
  # Add profile with explicit age public key
  sopsy profile add dev --description "Development" --age "age1..."
  
  # Add profile with multiple age recipients
  sopsy profile add team --age "age1abc..." --age "age1def..."`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		description, _ := cmd.Flags().GetString("description")
		ageKeys, _ := cmd.Flags().GetStringSlice("age")
		ageKeyFile, _ := cmd.Flags().GetString("age-key-file")

		profile := &config.Profile{
			Name:        name,
			Description: description,
		}

		// Add age backend
		if ageKeyFile != "" || len(ageKeys) > 0 {
			profile.Age = &config.AgeConfig{
				KeyFile:    ageKeyFile,
				Recipients: ageKeys,
			}
		}

		// Validate
		if !profile.HasBackends() {
			return fmt.Errorf("at least one encryption backend is required (--age-key-file or --age)")
		}

		// Add to config
		if err := cfg.AddProfile(profile); err != nil {
			return err
		}

		// Save
		path, _ := config.DefaultConfigPath()
		if cfgFile != "" {
			path = cfgFile
		}
		if err := config.Save(cfg, path); err != nil {
			return err
		}

		fmt.Printf("Profile '%s' added successfully\n", name)
		return nil
	},
}

var profileLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all profiles",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles := cfg.ListProfiles()
		if len(profiles) == 0 {
			fmt.Println("No profiles configured")
			return nil
		}

		// Sort by name
		sort.Slice(profiles, func(i, j int) bool {
			return profiles[i].Name < profiles[j].Name
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "NAME\tDESCRIPTION\tBACKENDS")
		for _, p := range profiles {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.Description, p.GetBackendSummary())
		}
		_ = w.Flush()

		return nil
	},
}

var profileShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show profile details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := cfg.GetProfile(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Name:        %s\n", profile.Name)
		fmt.Printf("Description: %s\n", profile.Description)
		fmt.Printf("Backends:    %s\n", profile.GetBackendSummary())

		if profile.Age != nil && len(profile.Age.Recipients) > 0 {
			fmt.Println("\nAge Recipients:")
			for _, r := range profile.Age.Recipients {
				fmt.Printf("  - %s\n", r)
			}
		}

		return nil
	},
}

var profileRmCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remove a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if err := cfg.RemoveProfile(name); err != nil {
			return err
		}

		path, _ := config.DefaultConfigPath()
		if cfgFile != "" {
			path = cfgFile
		}
		if err := config.Save(cfg, path); err != nil {
			return err
		}

		fmt.Printf("Profile '%s' removed\n", name)
		return nil
	},
}

var profileEditCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit a profile in your editor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := config.DefaultConfigPath()
		if cfgFile != "" {
			path = cfgFile
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}

		c := exec.Command(editor, path)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

var profileUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set default profile and output env export",
	Long: `Set a profile as the default and print environment export statement.

If no profile name is given, opens fzf to select interactively.

Examples:
  sopsy profile use stg       # Use specific profile
  sopsy profile use           # Select with fzf
   
After this, you can use standard SOPS commands:
  sops -e -i secrets.yaml
  sops -d secrets.yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string

		if len(args) == 0 {
			// Interactive selection with fzf
			profiles := cfg.ListProfiles()
			if len(profiles) == 0 {
				return fmt.Errorf("no profiles configured, run: sopsy profile add <name> --age-key-file <path>")
			}

			// Build list of profile names
			var names []string
			for _, p := range profiles {
				names = append(names, p.Name)
			}

			// Try fzf
			selected, err := selectWithFzf(names)
			if err != nil {
				// Fallback: list profiles
				fmt.Fprintln(os.Stderr, "Available profiles:")
				for _, n := range names {
					fmt.Fprintf(os.Stderr, "  - %s\n", n)
				}
				return fmt.Errorf("specify profile name: sopsy profile use <name>")
			}
			name = selected
		} else {
			name = args[0]
		}

		// Get profile (also verifies it exists)
		profile, err := cfg.GetProfile(name)
		if err != nil {
			return err
		}

		cfg.DefaultProfile = name

		path, _ := config.DefaultConfigPath()
		if cfgFile != "" {
			path = cfgFile
		}
		if err := config.Save(cfg, path); err != nil {
			return err
		}

		// Output export statement for shell integration
		if profile.Age != nil && profile.Age.KeyFile != "" {
			keyPath := profile.Age.GetKeyFilePath()
			fmt.Printf("export SOPS_AGE_KEY_FILE=\"%s\"\n", keyPath)
		}

		return nil
	},
}

// selectWithFzf uses fzf to select from a list of options
func selectWithFzf(options []string) (string, error) {
	// Check if fzf is available
	_, err := exec.LookPath("fzf")
	if err != nil {
		return "", fmt.Errorf("fzf not found")
	}

	cmd := exec.Command("fzf", "--height=10", "--prompt=Select profile: ")
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		for _, opt := range options {
			if _, err := fmt.Fprintln(stdin, opt); err != nil {
				return
			}
		}
		_ = stdin.Close()
	}()

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	selected := string(out)
	if len(selected) > 0 && selected[len(selected)-1] == '\n' {
		selected = selected[:len(selected)-1]
	}

	return selected, nil
}

var profileCurrentCmd = &cobra.Command{
	Use:    "current",
	Short:  "Show current default profile env",
	Hidden: true, // Internal use for shell integration
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.DefaultProfile == "" {
			return nil
		}

		profile, err := cfg.GetProfile(cfg.DefaultProfile)
		if err != nil {
			return nil // Silently fail if profile not found
		}

		if profile.Age != nil && profile.Age.KeyFile != "" {
			keyPath := profile.Age.GetKeyFilePath()
			fmt.Printf("export SOPS_AGE_KEY_FILE=\"%s\"\n", keyPath)
		}

		return nil
	},
}

var profileResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Clear default profile",
	Long:  `Clear the default profile. After reset, fzf will prompt for profile selection.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.DefaultProfile = ""

		path, _ := config.DefaultConfigPath()
		if cfgFile != "" {
			path = cfgFile
		}
		if err := config.Save(cfg, path); err != nil {
			return err
		}

		fmt.Println("Default profile cleared")
		return nil
	},
}

func init() {
	profileAddCmd.Flags().String("description", "", "profile description")
	profileAddCmd.Flags().String("age-key-file", "", "path to age key file (contains public and private keys)")
	profileAddCmd.Flags().StringSlice("age", nil, "age recipient public keys")

	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileLsCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileRmCmd)
	profileCmd.AddCommand(profileEditCmd)
	profileCmd.AddCommand(profileUseCmd)
	profileCmd.AddCommand(profileResetCmd)
	profileCmd.AddCommand(profileCurrentCmd)
}
