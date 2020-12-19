package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.com/schoentoon/rs-tools/lib/info"
)

var voiceOfSerenCmd = &cobra.Command{
	Use:   "voiceofseren",
	Short: "Retrieve the current voice of seren",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := info.VoiceOfSeren(http.DefaultClient)
		if err != nil {
			return err
		}
		fmt.Printf("%s, %s\n", res[0], res[1])
		return nil
	},
}
