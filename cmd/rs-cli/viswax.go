package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/info"
)

var viswaxCmd = &cobra.Command{
	Use:   "viswax",
	Short: "Retrieve the viswax combination for today",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := info.Viswax(http.DefaultClient)
		if err != nil {
			return err
		}
		fmt.Printf("Combination for %s\n", res.Date.Format("January 2 2006"))
		fmt.Printf("Primary rune is %s, costing %dgp\n", res.Primary.Rune, res.Primary.Cost)

		for i, rne := range res.Secondary {
			fmt.Printf("Secondary rune is either %s, costing %dgp", rne.Rune, rne.Cost)
			if i != (len(res.Secondary) - 1) {
				fmt.Printf(" or\n")
			} else {
				fmt.Printf("\n")
			}
		}

		fmt.Printf("Last rune is always random and unique per player, use your runecrafting cape to find out.\n")
		return nil
	},
}
