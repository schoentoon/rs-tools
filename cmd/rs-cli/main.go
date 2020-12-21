package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "rs-cli"}
	rootCmd.AddCommand(completionCmd)

	rootCmd.AddCommand(alogCmd)
	rootCmd.AddCommand(araxxorCmd)
	rootCmd.AddCommand(itemDBCmd)
	rootCmd.AddCommand(voiceOfSerenCmd)
	rootCmd.AddCommand(voragoCmd)
	rootCmd.AddCommand(spotlightCmd)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
