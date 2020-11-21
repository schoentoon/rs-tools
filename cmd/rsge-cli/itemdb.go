package main

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/c-bata/go-prompt"
	"gitlab.com/schoentoon/rs-tools/lib/ge/itemdb"
)

type ItemDB struct {
}

func (d *ItemDB) Name() string { return "download-itemdb" }

func (d *ItemDB) Description() string { return "Download and fill a local copy of the item database" }

func (d *ItemDB) Autocomplete(app *Application, in prompt.Document) []prompt.Suggest { return nil }

func (d *ItemDB) WantSpinner() bool { return false }

func (d *ItemDB) Execute(app *Application, argv string, out io.Writer) error {
	if argv == "" {
		return errors.New("You HAVE to provide a filename to save the result to")
	}

	f, err := os.OpenFile(argv, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	db := itemdb.New()
	db.Writer = f

	return db.Update(http.DefaultClient, 8)
}
