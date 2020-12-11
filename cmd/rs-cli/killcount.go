package main

import (
	"fmt"
	"io"
	"os"

	"github.com/adrg/xdg"
	"github.com/c-bata/go-prompt"
	"gitlab.com/schoentoon/rs-tools/lib/runemetrics"
)

type Killcount struct {
	KC map[string]int
}

func (k *Killcount) Name() string { return "killcount" }

func (k *Killcount) Description() string { return "Calculate the killcounts based on stored alogs" }

func (k *Killcount) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest {
	return nil
}

func (k *Killcount) WantSpinner() bool { return true }

func (k *Killcount) HandleActivity(activity runemetrics.Activity) error {
	kill, err := activity.BossKills()
	if err != nil { // error really just means it's probably not a bosskill, so we ignore it
		return nil
	}
	k.KC[kill.Boss] += kill.Amount
	return nil
}

func (k *Killcount) Execute(app *Application, argv string, out io.Writer) error {
	username := argv
	if username == "" || len(username) > 12 {
		return fmt.Errorf("You need to specify a valid username")
	}

	filename, err := xdg.DataFile(fmt.Sprintf("rscli/alog/%s.ljson", username))
	if err != nil {
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	k.KC = make(map[string]int)

	err = runemetrics.IterateActivities(f, k)
	if err != nil {
		return err
	}

	for boss, amount := range k.KC {
		fmt.Fprintf(out, "%s: %d\n", boss, amount)
	}

	return nil
}
