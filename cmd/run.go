package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts running the PintaDB server",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement
	},
}
