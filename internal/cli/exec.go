package cli

import (
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec -- <sops args...>",
	Short: "Execute SOPS with profile and custom arguments",
	Long: `Execute SOPS with the selected profile's configuration plus custom arguments.

Example:
  sopsctl exec -p myprofile -- --in-place secret.yaml`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := selectProfile()
		if err != nil {
			return err
		}

		baseArgs := make([]string, 0, 16)

		if profile.Age != nil {
			for _, recipient := range profile.Age.Recipients {
				baseArgs = append(baseArgs, "--age", recipient)
			}
		}

		sopsArgs := append(baseArgs, args...)

		if dryRun {
			executor.DryRun(sopsArgs)
			return nil
		}

		return executor.Execute(sopsArgs)
	},
}
