package main

import (
	"net/http"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
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

		dbf, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		defer dbf.Close()

		return meta.Download(http.DefaultClient, db, dbf)
	},
}
