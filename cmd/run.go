package cmd

import (
	"fmt"

	"github.com/columbusearch/pintadb/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts running the PintaDB server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting PintaDB v%s server...\n", server.Version)
		_, err := server.NewServer(cfg)
		if err != nil {
			// TODO replace with logger
			fmt.Println(err)
		}
	},
}
