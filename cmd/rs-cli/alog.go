package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"gitlab.com/schoentoon/rs-tools/lib/runemetrics"
)

type Alog struct {
}

func (a *Alog) Name() string { return "alog" }

func (a *Alog) Description() string { return "Retrieve the adventure log of a specified user" }

func (a *Alog) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest {
	path, err := xdg.DataFile("rscli/alog")
	if err != nil {
		return nil
	}
	files, err := filepath.Glob(fmt.Sprintf("%s/*.ljson", path))
	if err != nil {
		return nil
	}

	w := in.GetWordBeforeCursor()
	out := make([]prompt.Suggest, 0, len(files))
	for _, file := range files {
		file = strings.TrimSuffix(file, ".ljson")
		file = filepath.Base(file)
		out = append(out, prompt.Suggest{Text: file})
	}

	return prompt.FilterHasPrefix(out, w, true)
}

func (a *Alog) WantSpinner() bool { return true }

func (a *Alog) Execute(app *Application, argv string, out io.Writer) error {
	username := argv
	if username == "" || len(username) > 12 {
		return fmt.Errorf("You need to specify a valid username")
	}

	profile, err := runemetrics.FetchProfile(app.Client, username)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(out)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"When", "Activity"})
	defer table.Render()

	for _, activity := range profile.Activities {
		table.Append([]string{activity.Date.Local().Format("02-Jan-2006 15:04"), activity.Details})
	}

	return nil
}
