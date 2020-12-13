package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adrg/xdg"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"gitlab.com/schoentoon/rs-tools/lib/runemetrics"
)

type UpdateAlog struct {
}

func (a *UpdateAlog) Name() string { return "update-alog" }

func (a *UpdateAlog) Description() string {
	return "Updates the adventure log of a specified user in a local copy"
}

func (a *UpdateAlog) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest {
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

func (a *UpdateAlog) WantSpinner() bool { return true }

func (a *UpdateAlog) Execute(app *Application, argv string, out io.Writer) error {
	username := argv
	if username == "" || len(username) > 12 {
		return fmt.Errorf("You need to specify a valid username")
	}

	filename, err := xdg.DataFile(fmt.Sprintf("rscli/alog/%s.ljson", username))
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Base(filename), 0600)
	if err != nil {
		return err
	}

	profile, err := runemetrics.FetchProfile(app.Client, username)
	if err != nil {
		return err
	}

	existing, err := a.readOutputFile(filename)
	if err != nil {
		return err
	}

	newer := profile.Activities
	if len(existing) > 0 {
		newer = runemetrics.NewAchievementsSince(existing, profile.Activities)
		if len(newer) >= 20 {
			color.New(color.FgRed).Fprintf(out, "20 new activities, likely missing some in between!\n")
		}
	}

	fout, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer fout.Close()

	sort.Slice(newer, func(i, j int) bool { return newer[i].Date.Unix() < newer[j].Date.Unix() })

	return runemetrics.WriteActivities(fout, newer)
}

func (a *UpdateAlog) readOutputFile(filename string) ([]runemetrics.Activity, error) {
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