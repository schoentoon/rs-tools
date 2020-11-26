package main

import (
	"fmt"
	"io"

	"github.com/c-bata/go-prompt"
	"gitlab.com/schoentoon/rs-tools/lib/info"
)

type Araxxor struct {
}

func (a *Araxxor) Name() string { return "araxxor" }

func (a *Araxxor) Description() string { return "Retrieve the current open paths of araxxor" }

func (a *Araxxor) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest {
	return nil
}

func (a *Araxxor) WantSpinner() bool { return true }

func (a *Araxxor) Execute(app *Application, argv string, out io.Writer) error {
	res, err := info.AraxxorPath(app.Client)
	if err != nil {
		return err
	}

	openOrClosed := func(b bool) string {
		if b {
			return "open"
		}
		return "closed"
	}

	fmt.Fprintf(out, "%s\n", res.Description)
	fmt.Fprintf(out, "Minions path is %s\n", openOrClosed(res.Minions))
	fmt.Fprintf(out, "Acid path is %s\n", openOrClosed(res.Acid))
	fmt.Fprintf(out, "Darkness path is %s\n", openOrClosed(res.Darkness))
	fmt.Fprintf(out, "This rotation will change in %d days.\n", res.DaysLeft)

	return nil
}
