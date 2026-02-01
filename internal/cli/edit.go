package cli

import (
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <file>",
	Short: "Edit an encrypted file",
	Long: `Edit an encrypted file using SOPS.

Opens the file in your $EDITOR after decrypting, then re-encrypts on save.
If the file is new or being re-keyed, a profile can be selected.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		// For edit, profile is optional (SOPS can work with existing files)
		var sopsArgs []string
		if profileName != "" {
			profile, err := cfg.GetProfile(profileName)
			if err != nil {
				return err
			}
			sopsArgs = builder.BuildEdit(profile, file)
		} else {
			sopsArgs = []string{"edit", file}
		}

		// Execute or dry-run
		if dryRun {
			executor.DryRun(sopsArgs)
			return nil
		}

		return executor.Execute(sopsArgs)
	},
}
