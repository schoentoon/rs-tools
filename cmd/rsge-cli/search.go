package main

import (
	"fmt"
	"io"

	"github.com/c-bata/go-prompt"
)

type Search struct {
}

func (s *Search) Name() string { return "search" }

func (s *Search) Description() string { return "Search for an item in the GE database" }

func (s *Search) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest { return nil }

func (s *Search) WantSpinner() bool { return true }

func (s *Search) Execute(app *Application, argv string, out io.Writer) error {
	items, err := app.Ge.SearchItems(argv)
	if err != nil {
		return err
	}

	for _, item := range items {
		app.ItemCache[item.ItemID] = item.Name
		fmt.Fprintf(out, "%s - %d\n", item.Name, item.ItemID)
	}

	return nil
}
