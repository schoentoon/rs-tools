package main

import (
	"fmt"
	"io"

	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
)

type Search struct {
}

func (s *Search) Name() string { return "search" }

func (s *Search) Description() string { return "Search for an item in the GE database" }

func (s *Search) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest { return nil }

func (s *Search) WantSpinner() bool { return true }

func (s *Search) Execute(app *Application, argv string, out io.Writer) error {
	items, err := app.Search.SearchItems(argv)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"ID", "Item"})
	defer table.Render()

	for _, item := range items {
		app.ItemCache[item.ID] = item.Name
		table.Append([]string{fmt.Sprintf("%d", item.ID), item.Name})
	}

	return nil
}
