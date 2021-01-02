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

var alogClueScrolls = &cobra.Command{
	Use:   "cluescrolls",
	Short: "Calculate amount of clue scrolls done",

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

		c := &ClueScroll{
			Difficulties: make(map[runemetrics.ClueDifficulty]int),
		}

		err = runemetrics.IterateActivities(f, c)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Difficulty", "Amount"})
		defer table.Render()

		for difficulty, amount := range c.Difficulties {
			table.Append([]string{difficulty.String(), fmt.Sprintf("%d", amount)})
		}

		return nil
	},
}

type ClueScroll struct {
	Difficulties map[runemetrics.ClueDifficulty]int
}

func (c *ClueScroll) HandleActivity(activity runemetrics.Activity) error {
	clue, err := activity.ClueScroll()
	if err != nil { // error really just means it's probably not a bosskill, so we ignore it
		return nil
	}
	c.Difficulties[clue.Difficulty]++
	return nil
}
