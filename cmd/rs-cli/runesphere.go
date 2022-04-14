package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/info"
)

var runesphereCmd = &cobra.Command{
	Use:   "runesphere",
	Short: "Retrieve information about when the runesphere is going to be active next.",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := info.Runesphere(http.DefaultClient)
		if err != nil {
			return err
		}

		if res.Active {
			fmt.Printf("Runesphere is currently active\n")
		} else {
			fmt.Printf("The next Runesphere will be active in %s at %s\n", time.Until(res.Next), res.Next.Local())
		}

		return nil
	},
}
