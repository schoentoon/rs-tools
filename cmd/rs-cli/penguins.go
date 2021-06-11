package main

import (
	"net/http"
	"os"

	"github.com/ararog/timeago"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/info"
)

var penguinsCmd = &cobra.Command{
	Use:   "penguins",
	Short: "Retrieve the current penguins locations",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := info.Penguins(http.DefaultClient)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Disguise", "Last location", "Warning", "Last seen"})
		defer table.Render()

		for _, penguin := range res.ActivePenguins {
			last_seen, err := timeago.TimeAgoFromNowWithTime(penguin.LastSeen.Time)
			if err != nil {
				return err
			}
			table.Append([]string{penguin.Name, penguin.Disguise, penguin.LastLocation, penguin.Warning, last_seen})
		}

		for _, bear := range res.Bear {
			table.Append([]string{bear.Name, "Well", bear.Location, "", ""})
		}

		return nil
	},
}
