package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/runemetrics"
)

var alogCmd = &cobra.Command{
	Use:   "alog",
	Short: "Retrieve the adventure log of a specified user",

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Need at least a username")
		}
		username := args[0]
		if username == "" || len(username) > 12 {
			return fmt.Errorf("You need to specify a valid username")
		}
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		path, err := xdg.DataFile("rscli/alog")
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		files, err := filepath.Glob(fmt.Sprintf("%s/*.ljson", path))
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		out := make([]string, 0, len(files))
		for _, file := range files {
			file = strings.TrimSuffix(file, ".ljson")
			file = filepath.Base(file)
			if strings.HasPrefix(file, toComplete) {
				out = append(out, file)
			}
		}

		return out, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		profile, err := runemetrics.FetchProfile(http.DefaultClient, username)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoWrapText(false)
		table.SetHeader([]string{"When", "Activity"})
		defer table.Render()

		for _, activity := range profile.Activities {
			table.Append([]string{activity.Date.Local().Format("02-Jan-2006 15:04"), activity.Details})
		}
		return nil
	},
}

func init() {
	alogCmd.AddCommand(alogUpdate)
	alogCmd.AddCommand(alogKillCount)
	alogCmd.AddCommand(alogClueScrolls)
}
