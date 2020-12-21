package main

import (
	"github.com/spf13/cobra"
)

var spotlightCmd = &cobra.Command{
	Use:   "spotlight",
	Short: "Check the current and future spotlights",
}

func init() {
	spotlightCmd.AddCommand(spotlightMinigamesCmd)
}
