/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	api "mox-monoline/cmd/mox-monolith"
)

// cardNameCmd represents the cardName command
var cardNameCmd = &cobra.Command{
	Use:   "cardName",
	Short: "Search for a list of cards based on the search text",
	Long:  `e.g. mox cardName "Sol RIng"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("name argument missing, must not be blank.")
			return
		}

		if len(args) > 0 && args[0] == "" {
			fmt.Println("name argument missing, must not be blank.")
			return
		}

		cardNameData := api.FuzzySearchCardNames(args[0])

		for i := 0; i < len(cardNameData.Cards); i++ {
			fmt.Println(cardNameData.Cards[i].Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(cardNameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cardNameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cardNameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
