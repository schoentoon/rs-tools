package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/c-bata/go-prompt"
	"gitlab.com/schoentoon/rs-tools/lib/ge"
)

type Search struct {
}

func (s *Search) Name() string { return "search" }

func (s *Search) Description() string { return "Search for an item in the GE database" }

func (s *Search) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest { return nil }

func (s *Search) WantSpinner() bool { return true }

func (s *Search) Execute(app *Application, argv string, out io.Writer) error {
	items, err := ge.SearchItems(argv, http.DefaultClient)
	if err != nil {
		return err
	}

	for _, item := range items {
		app.ItemCache[item.ItemID] = item.Name
		fmt.Fprintf(out, "%s - %d\n", item.Name, item.ItemID)
	}

	return nil
}
