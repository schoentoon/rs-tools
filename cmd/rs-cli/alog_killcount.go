package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/runemetrics"
)

var alogKillCount = &cobra.Command{
	Use:   "killcount",
	Short: "Calculate the killcounts based on stored alogs",

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

		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		k := &Killcount{
			KC: make(map[string][]int),
		}

		err = runemetrics.IterateActivities(f, k)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Boss", "Kills", "HM"})
		defer table.Render()

		for boss, kills := range k.KC {
			table.Append([]string{boss, fmt.Sprintf("%d", kills[0]), fmt.Sprintf("%d", kills[1])})
		}

		return nil
	},
}

type Killcount struct {
	// the string will be the name of the boss, the []int will always only have 2 entries.
	// the first entry will be amount of kills in normal mode, second is amount of kills in hard mode
	// challenge mode whatever they call it for a specific boss, for this reason it could also be
	// completely ignored
	KC map[string][]int
}

func (k *Killcount) HandleActivity(activity runemetrics.Activity) error {
	kill, err := activity.BossKills()
	if err != nil { // error really just means it's probably not a bosskill, so we ignore it
		return nil
	}
	if k.KC[kill.Boss] == nil {
		k.KC[kill.Boss] = []int{0, 0}
	}
	if kill.Hardmode {
		k.KC[kill.Boss][1] += kill.Amount
	} else {
		k.KC[kill.Boss][0] += kill.Amount
	}
	return nil
}
