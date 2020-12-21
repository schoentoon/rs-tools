package main

import (
	"net/http"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/spotlight"
)

var spotlightMinigamesCmd = &cobra.Command{
	Use:   "minigames",
	Short: "Retrieve the current and future minigame spotlights",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := spotlight.Minigames(http.DefaultClient)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Minigame", "When"})
		defer table.Render()

		table.Append([]string{res.Current, "Now"})
		return res.Iterate(func(when time.Time, minigame string) error {
			table.Append([]string{minigame, when.Format("Mon Jan 2 2006")})
			return nil
		})
	},
}
