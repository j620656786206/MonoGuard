package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize MonoGuard configuration",
	Long: `Create a .monoguard.yaml configuration file in the
current directory with default settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		format := viper.GetString("format")
		if format == "json" {
			fmt.Fprintln(cmd.OutOrStdout(), `{"status":"placeholder","message":"Init will be implemented in Epic 8"}`)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "🚀 MonoGuard Init (Placeholder)\n")
			fmt.Fprintf(cmd.OutOrStdout(), "   Status: Will be implemented in Epic 8\n")
		}
	},
}

// Command registration is handled by root.go registerCommands()
