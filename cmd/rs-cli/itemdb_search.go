package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var itemDBSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for an item in the itemdb",

	RunE: func(cmd *cobra.Command, args []string) error {
		itemDB := readItemDB()

		items, err := itemDB.SearchItems(strings.Join(args, " "))
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Item"})
		defer table.Render()

		for _, item := range items {
			table.Append([]string{fmt.Sprintf("%d", item.ID), item.Name})
		}

		return nil
	},
}
