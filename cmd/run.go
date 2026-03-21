package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "To run the db-sync",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting sync...")
	},
}
