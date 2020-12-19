package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/info"
)

var voragoCmd = &cobra.Command{
	Use:   "vorago",
	Short: "Retrieve the current rotation of vorago",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := info.VoragoRotation(http.DefaultClient)
		if err != nil {
			return err
		}

		fmt.Printf("Current rotation is %s\n", res.Rotation)
		fmt.Printf("This rotation will change in %d days\n", res.DaysLeft)

		return nil
	},
}
