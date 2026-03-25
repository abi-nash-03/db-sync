package cmd

import (
	"db-sync/config"
	"fmt"
	"github.com/spf13/cobra"
)

var version = "dev"
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
	rootCmd.Version = version
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to the config file")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run the command")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func run(cmd *cobra.Command, args []string) {
	err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	if dryRun {
		fmt.Println("Dry run")
		if configPath != "" {
			fmt.Printf("[dry-run] Would load config from: %s\n", configPath)
			fmt.Println("[dry-run] No changes will be made.")
		}
		return
	}
}
