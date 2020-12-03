package main

import (
	"fmt"
	"io"

	"github.com/c-bata/go-prompt"
	"gitlab.com/schoentoon/rs-tools/lib/info"
)

type Vorago struct {
}

func (v *Vorago) Name() string { return "vorago" }

func (v *Vorago) Description() string { return "Retrieve the current rotation of vorago" }

func (v *Vorago) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest {
	return nil
}

func (v *Vorago) WantSpinner() bool { return true }

func (v *Vorago) Execute(app *Application, argv string, out io.Writer) error {
	res, err := info.VoragoRotation(app.Client)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Current rotation is %s\n", res.Rotation)
	fmt.Fprintf(out, "This rotation will change in %d days\n", res.DaysLeft)

	return nil
}
