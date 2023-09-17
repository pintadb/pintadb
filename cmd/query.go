package cmd

import (
	"fmt"

	"github.com/columbusearch/pintadb/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().IntVarP(&k, "maxResults", "", 10, "Number of results to return")

}

var (
	k int
)

// TODO - Make this a client command
var queryCmd = &cobra.Command{
	Use:   "query [string of query]",
	Short: "Query the PintaDB database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO replace with client
		pinta, err := server.NewServer(cfg)
		if err != nil {
			fmt.Println(err)
		}

		results, err := pinta.Search(args[0], k)
		if err != nil {
			panic(err)
		}

		// Print results
		for i, r := range results {
			if i == 0 {
				fmt.Printf("[%d] %s\n", r, pinta.Documents[r].RawText)
			} else {
				fmt.Printf("[%d]\n", r)
			}
		}
	},
}
