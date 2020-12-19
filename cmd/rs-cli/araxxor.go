package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/info"
)

var araxxorCmd = &cobra.Command{
	Use:   "araxxor",
	Short: "Retrieve the current open paths of araxxor",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := info.AraxxorPath(http.DefaultClient)
		if err != nil {
			return err
		}

		openOrClosed := func(b bool) string {
			if b {
				return "open"
			}
			return "closed"
		}

		fmt.Printf("%s\n", res.Description)
		fmt.Printf("Minions path is %s\n", openOrClosed(res.Minions))
		fmt.Printf("Acid path is %s\n", openOrClosed(res.Acid))
		fmt.Printf("Darkness path is %s\n", openOrClosed(res.Darkness))
		fmt.Printf("This rotation will change in %d days.\n", res.DaysLeft)

		return nil
	},
}
