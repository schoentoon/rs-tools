package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/ge"
	"gitlab.com/schoentoon/rs-tools/lib/ge/download"
)

const ITEMDB_LOCATION = "rscli/itemdb.ljson"

func readItemDB() ge.SearchItemInterface {
	filename, err := xdg.DataFile(ITEMDB_LOCATION)
	if err != nil {
		return geApi()
	}
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return geApi()
	}
	defer f.Close()

	out, err := download.NewDBFromReader(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return geApi()
	}
	return out
}

func geApi() *ge.Ge {
	return &ge.Ge{
		Client: http.DefaultClient,
		// It's not very nice to 'abuse' the firefox user agent here.. but for the only not really api
		// call they have on the ge website a captcha tended to get in the way sometimes. on first sight
		// switching to this user agent seemed to work around it, nasty but it works I guess
		// just don't call Search too often because of this really
		UserAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	}
}

var itemDBCmd = &cobra.Command{
	Use:   "itemdb",
	Short: "Interact with a local copy of the item database",
}

func init() {
	itemDBCmd.AddCommand(itemDBDownloadCmd)
	itemDBCmd.AddCommand(itemDBSearchCmd)
	itemDBCmd.AddCommand(itemDBPriceCmd)
}
