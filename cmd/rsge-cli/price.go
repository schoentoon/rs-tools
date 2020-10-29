package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"gitlab.com/schoentoon/rs-tools/lib/ge"
)

type Price struct {
}

func (p *Price) Name() string { return "price" }

func (p *Price) Description() string { return "Retrieve the current price of specified item" }

func (p *Price) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	out := make([]prompt.Suggest, 0, len(app.ItemCache))
	for _, name := range app.ItemCache {
		out = append(out, prompt.Suggest{Text: name})
	}

	return prompt.FilterHasPrefix(out, w, true)
}

func (p *Price) WantSpinner() bool { return true }

func (p *Price) Execute(app *Application, argv string, out io.Writer) error {
	id, err := strconv.ParseInt(argv, 10, 64)
	if err != nil {
		// as we're clearly not an item id we'll go look through our item cache
		// maybe we have an item with this name
		id = -1
		lower := strings.ToLower(argv)
		for i, name := range app.ItemCache {
			if strings.ToLower(name) == lower {
				id = i
				break
			}
		}
		// if it still doesn't look like it we go do a lookup
		if id == -1 {
			search, err := ge.SearchItems(argv, http.DefaultClient)
			if err != nil {
				return err
			}

			for _, item := range search {
				if id == -1 {
					id = item.ItemID
				}
				app.ItemCache[item.ItemID] = item.Name
			}
		}
		if id == -1 {
			return fmt.Errorf("No item found")
		}
	}

	graph, err := ge.PriceGraph(id, http.DefaultClient)
	if err != nil {
		return err
	}

	_, latest := graph.LatestPrice()
	fmt.Fprintf(out, "%d\n", latest)

	return nil
}
