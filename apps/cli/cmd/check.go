package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkOutput represents the JSON output structure for check command
type checkOutput struct {
	Status  string `json:"status"`
	Path    string `json:"path"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

var (
	failOn    string
	threshold int
)

var checkCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Validate dependencies for CI/CD",
	Long: `Run validation checks on the monorepo dependencies.
Returns exit code 0 on success, 1 on failure.
Designed for CI/CD integration.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		format := viper.GetString("format")
		if format == "json" {
			output := checkOutput{
				Status:  "placeholder",
				Path:    path,
				Passed:  true,
				Message: "Check will be implemented in Epic 2",
			}
			jsonBytes, _ := json.Marshal(output)
			fmt.Fprintln(cmd.OutOrStdout(), string(jsonBytes))
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "✅ MonoGuard Check (Placeholder)\n")
			fmt.Fprintf(cmd.OutOrStdout(), "   Path: %s\n", path)
			fmt.Fprintf(cmd.OutOrStdout(), "   Status: Passed (placeholder)\n")
		}
		// Exit code 0 is default when no error is returned
	},
}

func init() {
	// Command registration is handled by root.go registerCommands()
	// Local flags are registered here
	checkCmd.Flags().StringVar(&failOn, "fail-on", "all",
		"fail on: circular|boundary|all")
	checkCmd.Flags().IntVar(&threshold, "threshold", 0,
		"fail if health score below threshold (0-100)")
}
