package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"gitlab.com/schoentoon/rs-tools/lib/runemetrics"
)

type Killcount struct {
	// the string will be the name of the boss, the []int will always only have 2 entries.
	// the first entry will be amount of kills in normal mode, second is amount of kills in hard mode
	// challenge mode whatever they call it for a specific boss, for this reason it could also be
	// completely ignored
	KC map[string][]int
}

func (k *Killcount) Name() string { return "killcount" }

func (k *Killcount) Description() string { return "Calculate the killcounts based on stored alogs" }

func (k *Killcount) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest {
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

func (k *Killcount) WantSpinner() bool { return true }

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

	k.KC = make(map[string][]int)

	err = runemetrics.IterateActivities(f, k)
	if err != nil {
		return err
	}

	for boss, kills := range k.KC {
		fmt.Fprintf(out, "%s: %d", boss, kills[0])
		if kills[1] > 0 {
			color.New(color.FgRed).Fprintf(out, " (%d)", kills[1])
		}
		fmt.Fprintf(out, "\n")
	}

	return nil
}
