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

var stores string

// deckCmd represents the deck command
var deckCmd = &cobra.Command{
	Use:   "deck",
	Short: "Fetches the cheapest available card for each card in a given decklist",
	Long:  `e.g. mox-monoline deck --stores=Topdeck "/path/to/decklist.txt"`,
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
			// var total int32 = 0
			fmt.Println("card,Store1,Price1,Store2,Price2,Store3,Price3")
			for i := 0; i < len(cards); i++ {
				cardNameData := api.FuzzySearchCardNames(cards[i])
				cheapestCard, secondCheapestCard, thirdCheapestCard := api.FindCheapestCardAtRetailers(cardNameData.Cards[0].Id, stores)
				var printoutCard string = cards[i]
				if cheapestCard != nil {
					printoutCard += ";" + cheapestCard.Retailer_name + ";" + fmt.Sprint(cheapestCard.Price/100.00)
					// total += cheapestCard.Price
				} else {
					printoutCard += ";;"
				}
				if secondCheapestCard != nil {
					printoutCard += ";" + secondCheapestCard.Retailer_name + ";" + fmt.Sprint(secondCheapestCard.Price/100.00)
				} else {
					printoutCard += ";;"
				}
				if thirdCheapestCard != nil {
					printoutCard += ";" + thirdCheapestCard.Retailer_name + ";" + fmt.Sprint(thirdCheapestCard.Price/100.00)
				} else {
					printoutCard += ";;"
				}
				fmt.Println(printoutCard)
			}
			// fmt.Printf("Total: R%d (+- ??? cards)", total/100)
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
	deckCmd.Flags().StringVar(&stores, "stores", "all", "comma seperated list of stores you wish to search e.g. --stores=TopDeck,Dracoti")
}
