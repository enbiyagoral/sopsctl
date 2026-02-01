package cli

import (
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt <file>",
	Short: "Decrypt a file",
	Long: `Decrypt a file using SOPS.

SOPS typically auto-detects the encryption method from the file metadata,
so a profile is usually not needed for decryption.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		// Build SOPS arguments (decrypt usually doesn't need profile)
		sopsArgs := builder.BuildDecrypt(file)

		// Execute or dry-run
		if dryRun {
			executor.DryRun(sopsArgs)
			return nil
		}

		return executor.Execute(sopsArgs)
	},
}
