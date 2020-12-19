package main

import (
	"net/http"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"gitlab.com/schoentoon/rs-tools/lib/ge/itemdb"
)

var itemDBDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download and fill a local copy of the item database",

	RunE: func(cmd *cobra.Command, args []string) error {
		filename, err := xdg.DataFile(ITEMDB_LOCATION)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		db := itemdb.New()
		db.Writer = f

		p := mpb.New(mpb.WithWidth(64), mpb.WithOutput(os.Stdout))
		tasks := p.AddBar(int64(0),
			mpb.PrependDecorators(
				decor.CountersNoUnit("%d / %d "),
				decor.Percentage(),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.AverageETA(decor.ET_STYLE_MMSS, decor.WC{W: 4}), "done",
				),
			),
		)

		progCh := make(chan *itemdb.Progress, 1)
		errCh := make(chan error, 1)

		go func() {
			errCh <- db.Update(http.DefaultClient, 1, progCh)
		}()

		for progress := range progCh {
			tasks.SetTotal(progress.Tasks, progress.Tasks == progress.Finished)
			tasks.SetCurrent(progress.Finished)
		}

		p.Wait()

		return <-errCh
	},
}
