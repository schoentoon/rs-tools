package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var itemDBPriceCmd = &cobra.Command{
	Use:   "price",
	Short: "Retrieve the current price of specified item",

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		itemDB := readItemDB()

		items, err := itemDB.SearchItems(strings.Join(args, " "))
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		out := make([]string, 0, len(items))
		for _, item := range items {
			out = append(out, item.Name)
		}

		return out, cobra.ShellCompDirectiveNoFileComp
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		itemDB := readItemDB()

		items, err := itemDB.SearchItems(strings.Join(args, " "))
		if err != nil {
			return err
		}

		ge := geApi()

		if len(items) == 1 {
			graph, err := ge.PriceGraph(items[0].ID)
			if err != nil {
				return err
			}

			_, latest := graph.LatestPrice()

			// TODO check whether we're running non interactive, be less pretty in that case
			prettyPrintPrice(os.Stdout, latest)
		}

		return nil
	},
}
