/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	api "mox-monoline/cmd/mox-monolith"
	"os"
	"strings"
)

// deckCmd represents the deck command
var deckCmd = &cobra.Command{
	Use:   "deck",
	Short: "Fetches the cheapest available card for each card in a given decklist",
	Long:  `e.g. mox-monoline deck "/path/to/decklist.txt"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing decklist argument. Must not be blank")
			return
		}

		if len(args) > 0 && args[0] == "" {
			fmt.Println("Missing decklist argumet. Must not be blank")
			return
		}

		data, err := os.ReadFile(args[0])

		if err != nil {
			fmt.Println("Failed to read given file. Closting.")
			os.Exit(1)
		}

		if data != nil {
			textContents := string(data)
			cards := strings.Split(textContents, "\r\n")
			var total int32 = 0
			for i := 0; i < len(cards); i++ {
				cardNameData := api.FuzzySearchCardNames(cards[i])
				cheapestCard := api.FindCheapestCardAtRetailers(cardNameData.Cards[0].Id)
				if cheapestCard != nil {
					fmt.Println(cards[i] + " - " + cheapestCard.Retailer_name + " - " + cheapestCard.PriceRead)
					total += cheapestCard.Price
				} else {
					fmt.Println(cards[i] + " - Card not available")
				}
			}
			fmt.Printf("Total: R%d (+- ??? cards)", total/100)
		}

	},
}

func init() {
	rootCmd.AddCommand(deckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
