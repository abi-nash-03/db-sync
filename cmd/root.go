package cmd

import (
	"db-sync/config"
	"db-sync/pipeline"
	"db-sync/scheduler"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var version string
var configPath string
var dryRun bool
var rootCmd = &cobra.Command{
	Use:   "db-sync [option]",
	Short: "A Powerful DB Sync CLI Application",
	Long: `DB Sync is a simple CLI tool to manage your database synchronization.

This can be used for single time syncronization or for continuous syncronization.`,
	Args: cobra.ArbitraryArgs,
	Run:  run,
}

func init() {
	// Persistent flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to the config file")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run the command")
}

// SetVersion sets the version string for the root command
func SetVersion(v string) {
	version = v
	rootCmd.Version = v
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func run(cmd *cobra.Command, args []string) {
	err := config.LoadConfig(configPath)

	if err != nil {
		slog.Error("Error loading config:", "error", err)
		os.Exit(1)

	}

	if config.AppConfig.Schedule != "" {
		if err := scheduler.Start(config.AppConfig, config.AppConfig.Schedule); err != nil {
			slog.Error("Error running scheduler:", "error", err)
			os.Exit(1)
		}
	} else {
		if err := pipeline.Run(config.AppConfig, dryRun); err != nil {
			slog.Error("Error running pipeline:", "error", err)
			os.Exit(1)
		}
	}
}
