package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/runemetrics"
)

var alogUpdate = &cobra.Command{
	Use:   "update",
	Short: "Updates the adventure log of a specified user in a local copy",

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

		filename, err := xdg.DataFile(fmt.Sprintf("rscli/alog/%s.ljson", username))
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(filename), 0600)
		if err != nil {
			return err
		}

		profile, err := runemetrics.FetchProfile(http.DefaultClient, username)
		if err != nil {
			return err
		}

		existing, err := readOutputFile(filename)
		if err != nil {
			return err
		}

		newer := profile.Activities
		if len(existing) > 0 {
			newer = runemetrics.NewAchievementsSince(existing, profile.Activities)
			if len(newer) >= 20 {
				_, _ = color.New(color.FgRed).Printf("20 new activities, likely missing some in between!\n")
			}
		}

		activities := runemetrics.Sort(existing, newer)

		fout, err := os.OpenFile(fmt.Sprintf("%s.tmp", filename), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
		if err != nil {
			return err
		}
		defer fout.Close()

		err = runemetrics.WriteActivities(fout, activities)
		if err != nil {
			return err
		}

		return os.Rename(fmt.Sprintf("%s.tmp", filename), filename)
	},
}

func readOutputFile(filename string) ([]runemetrics.Activity, error) {
	out := []runemetrics.Activity{}

	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	defer f.Close()

	return runemetrics.ReadActivities(f)
}
