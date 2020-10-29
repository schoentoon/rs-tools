package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
)

type Command interface {
	Name() string
	WantSpinner() bool
	Description() string
	Autocomplete(app *Application, in prompt.Document) []prompt.Suggest
	Execute(app *Application, args string, out io.Writer) error
}

type Application struct {
	Commands  []Command
	ItemCache map[int64]string
	Pretty    bool
}

func (a *Application) completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	line := in.TextBeforeCursor()
	blocks := strings.SplitN(line, " ", 2)
	cmd := blocks[0]

	for _, c := range a.Commands {
		if cmd == c.Name() {
			return c.Autocomplete(a, in)
		}
	}

	if len(blocks) != 1 {
		return []prompt.Suggest{}
	}

	out := make([]prompt.Suggest, len(a.Commands))
	for i, c := range a.Commands {
		out[i] = prompt.Suggest{Text: c.Name(), Description: c.Description()}
	}
	return prompt.FilterHasPrefix(out, w, true)
}

type StopSpinnerWriter struct {
	Out     io.Writer // this is the underlying writer
	Spinner *spinner.Spinner
}

func (w *StopSpinnerWriter) Write(p []byte) (int, error) {
	w.CloseSpinner()
	return w.Out.Write(p)
}

func (w *StopSpinnerWriter) CloseSpinner() {
	if w.Spinner != nil {
		w.Spinner.Stop()
		w.Spinner = nil
	}
}

func (a *Application) executor(in string) {
	in = strings.TrimSpace(in)

	blocks := strings.SplitN(in, " ", 2)
	if blocks[0] == "" {
		return
	}

	var cmd Command
	for _, c := range a.Commands {
		if c.Name() == blocks[0] {
			cmd = c
			break
		}
	}

	if cmd == nil {
		fmt.Println("Invalid command")
		return
	}

	out := &StopSpinnerWriter{Out: os.Stdout}
	if cmd.WantSpinner() && a.Pretty {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.HideCursor = true
		s.Start()
		out.Spinner = s
		defer out.CloseSpinner()
	}

	var err error
	if len(blocks) == 1 {
		err = cmd.Execute(a, "", out)
	} else {
		err = cmd.Execute(a, blocks[1], out)
	}
	if err == flag.ErrHelp {
		return
	}
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func main() {
	flag.Parse()
	a := Application{
		Commands: []Command{
			&Search{},
			&Price{},
		},
		ItemCache: make(map[int64]string),
		Pretty:    flag.NArg() == 0, // if we don't have any left over flags we're gonna be interactive
	}

	if flag.NArg() > 0 {
		a.executor(strings.Join(flag.Args(), " "))
		return
	}

	p := prompt.New(
		a.executor,
		a.completer,
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)

	p.Run()
}
