package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/adrg/xdg"
	"github.com/c-bata/go-prompt"
	"gitlab.com/schoentoon/rs-tools/lib/runemetrics"
)

type Alog struct {
}

func (a *Alog) Name() string { return "alog" }

func (a *Alog) Description() string { return "Update the adventure log of a certain user" }

func (a *Alog) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest {
	return nil
}

func (a *Alog) WantSpinner() bool { return true }

func (a *Alog) Execute(app *Application, argv string, out io.Writer) error {
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
	}

	fout, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer fout.Close()

	sort.Slice(newer, func(i, j int) bool { return newer[i].Date.Unix() < newer[j].Date.Unix() })

	return runemetrics.WriteActivities(fout, newer)
}

func (a *Alog) readOutputFile(filename string) ([]runemetrics.Activity, error) {
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
