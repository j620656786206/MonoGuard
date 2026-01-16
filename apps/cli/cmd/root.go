package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	verbose   bool
	format    string
	version   = "0.1.0" // Set via ldflags during build
	commit    = "unknown"
	buildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "monoguard",
	Short: "MonoGuard - Monorepo dependency analysis and validation",
	Long: `MonoGuard is a comprehensive tool for analyzing monorepo
dependencies, detecting circular dependencies, and providing
actionable fix suggestions.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	registerFlags()
	registerCommands()
}

// registerFlags adds all persistent flags to rootCmd
func registerFlags() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is .monoguard.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"verbose output")
	rootCmd.PersistentFlags().StringVar(&format, "format", "text",
		"output format (text|json)")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
}

// registerCommands adds all subcommands to rootCmd
func registerCommands() {
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(fixCmd)
	rootCmd.AddCommand(initCmd)
}

// ResetForTesting resets and re-registers all commands and flags
// This is used by tests that need a clean rootCmd state
func ResetForTesting() {
	rootCmd.ResetCommands()
	rootCmd.ResetFlags()

	// Reset global flag variables to defaults
	cfgFile = ""
	verbose = false
	format = "text"

	// Reset command-specific flags
	resetCheckFlags()
	resetFixFlags()

	registerFlags()
	registerCommands()
}

// resetCheckFlags resets check command flags to defaults
func resetCheckFlags() {
	failOn = "all"
	threshold = 0
}

// resetFixFlags resets fix command flags to defaults
func resetFixFlags() {
	dryRun = false
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory
		viper.AddConfigPath(".")
		// Search in home directory
		home, _ := os.UserHomeDir()
		viper.AddConfigPath(home + "/.monoguard")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".monoguard")
	}

	// Read environment variables with MONOGUARD_ prefix
	viper.SetEnvPrefix("MONOGUARD")
	viper.AutomaticEnv()

	// Read config file (ignore error if not found)
	viper.ReadInConfig()
}
