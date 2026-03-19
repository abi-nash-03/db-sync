package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:     "db-syn <option>",
	Short:   "DB Sync CLI Application",
	Long:    "DB Sync is a simple CLI tool to manage your database synchronization.",
	Example: "db-sync run",
	// Args: cobra.ExactArgs(1),
	// Run: func(cmd *cobra.Command, args []string) {
	//     option := args[0]
	//     switch option {
	//     case "run":
	//         fmt.Println("Starting sync...")
	//     case "validate":
	//         fmt.Println("Validating connections...")
	//     default:
	//         fmt.Println("invalid argument")
	//     }
	// },
}

func init() {
	rootCmd.AddCommand(runCmd)
}


func Execute() {
 cobra.CheckErr(rootCmd.Execute())
}