package main

import (
	"fmt"
	"io"

	"github.com/c-bata/go-prompt"
	"gitlab.com/schoentoon/rs-tools/lib/info"
)

type VoiceOfSeren struct {
}

func (v *VoiceOfSeren) Name() string { return "voiceofseren" }

func (v *VoiceOfSeren) Description() string { return "Retrieve the current voice of seren" }

func (v *VoiceOfSeren) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest {
	return nil
}

func (v *VoiceOfSeren) WantSpinner() bool { return true }

func (v *VoiceOfSeren) Execute(app *Application, argv string, out io.Writer) error {
	res, err := info.VoiceOfSeren(app.Client)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s, %s\n", res[0], res[1])
	return nil
}
