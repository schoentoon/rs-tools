package main

import (
	"net/http"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"gitlab.com/schoentoon/rs-tools/lib/ge/download"
)

const META_LOCATION = "rscli/item.meta.json"

func init() {
	filename, err := xdg.DataFile(ITEMDB_LOCATION)
	if err != nil {
		panic(err)
	}
	metafile, err := xdg.DataFile(META_LOCATION)
	if err != nil {
		panic(err)
	}
	itemDBDownloadCmd.PersistentFlags().StringP("meta", "m", metafile, "Location to write down the metadata file")
	itemDBDownloadCmd.PersistentFlags().StringP("file", "f", filename, "Location to write down the item database file")
}

var itemDBDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download and fill a local copy of the item database",

	RunE: func(cmd *cobra.Command, args []string) error {
		filename, err := cmd.PersistentFlags().GetString("file")
		if err != nil {
			return err
		}
		metafile, err := cmd.PersistentFlags().GetString("meta")
		if err != nil {
			return err
		}

		meta, err := download.DiffMetadataFromFile(http.DefaultClient, metafile)
		if err != nil {
			return err
		}

		metaf, err := os.OpenFile(metafile, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			return err
		}
		defer metaf.Close()

		err = meta.Serialize(metaf)
		if err != nil {
			return err
		}

		var db *download.DB
		dbf, err := os.Open(filename)
		if err != nil && !os.IsNotExist(err) {
			return err
		} else if dbf != nil {
			defer dbf.Close()
			db, err = download.NewDBFromReader(dbf)
			if err != nil {
				return err
			}
		} else {
			db = download.NewEmptyDB()
		}

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

		progCh := make(chan *download.Progress, 1)
		errCh := make(chan error, 1)

		dbf, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		defer dbf.Close()

		go func() {
			errCh <- meta.Download(http.DefaultClient, db, dbf, progCh)
		}()

		for progress := range progCh {
			tasks.SetTotal(progress.Tasks, progress.Tasks == progress.Finished)
			tasks.SetCurrent(progress.Finished)
		}

		p.Wait()

		return <-errCh
	},
}
