package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type CardNameData struct {
	Cards []Card
}

type Card struct {
	Id           int32
	Name         string
	Last_scraped string
	Is_scraping  bool
}

type CardPriceAvailabilityData struct {
	Products []MerchantCardUnit
	Card     Card
}

type MerchantCardUnit struct {
	Id            int32
	Name          string
	Price         int32
	PriceRead     string
	Link          string
	Stock         int32
	Is_foil       bool
	Last_scraped  string
	Image         string
	Retailer_id   int32
	Retailer_name string
}

func FuzzySearchCardNames(searchCriteria string) CardNameData {
	response, err := http.Get("https://moxmonolith.com/card/search?name=" + searchCriteria)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var cardNameData CardNameData

	json.Unmarshal([]byte(responseData), &cardNameData)

	return cardNameData
}

func constructStoresList(stores string) string {

	if stores == "all" {
		return "retailers^[^]=2&retailers^[^]=3&retailers^[^]=4&retailers^[^]=6&retailers^[^]=11&retailers^[^]=13&retailers^[^]=15&retailers^[^]=16&retailers^[^]=18&retailers^[^]=19&retailers^[^]=20&retailers^[^]=21&retailers^[^]=22&retailers^[^]=23&retailers^[^]=24&retailers^[^]=25&retailers^[^]=26&retailers^[^]=32&retailers^[^]=33&retailers^[^]=34&retailers^[^]=35"
	}

	availableStores := map[string]int{
		"TopDeck":               2,
		"Luckshack":             3,
		"TheWarren":             4,
		"BattleBunkerPaarl":     6,
		"Dracoti":               11,
		"D20Battleground":       13,
		"UnderworldConnections": 15,
		"TheStoneDragon":        16,
		"Sword&Board":           18,
		"MagicBazaar":           19,
		"TCGTrader":             20,
		"UntappedLands":         21,
		"GreedyGold":            26,
		"MirageGaming":          34,
		"GeekHome":              36,
	}

	retailersList := ""
	selectedRetailers := strings.Split(stores, ",")

	for i := 0; i < len(selectedRetailers); i++ {
		if i == 0 {
			retailersList += "retailers[]=" + fmt.Sprint(availableStores[selectedRetailers[i]])
		} else {
			retailersList += "&retailers[]=" + fmt.Sprint(availableStores[selectedRetailers[i]])
		}
	}

	return retailersList
}

func FindCheapestCardAtRetailers(cardId int32, stores string) (*MerchantCardUnit, *MerchantCardUnit, *MerchantCardUnit) {

	url := fmt.Sprint("https://moxmonolith.com/card/", cardId, "/products?"+constructStoresList(stores))

	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var cardPriceAvailabilityData CardPriceAvailabilityData

	json.Unmarshal([]byte(responseData), &cardPriceAvailabilityData)

	if len(cardPriceAvailabilityData.Products) == 1 {
		return &cardPriceAvailabilityData.Products[0], nil, nil
	}

	if len(cardPriceAvailabilityData.Products) == 2 {
		if cardPriceAvailabilityData.Products[0].Price < cardPriceAvailabilityData.Products[1].Price {
			return &cardPriceAvailabilityData.Products[0], &cardPriceAvailabilityData.Products[1], nil
		} else {
			return &cardPriceAvailabilityData.Products[1], &cardPriceAvailabilityData.Products[0], nil
		}
	}

	if len(cardPriceAvailabilityData.Products) > 0 {

		for i := 0; i < len(cardPriceAvailabilityData.Products); i++ {
			for j := 0; j < len(cardPriceAvailabilityData.Products); j++ {
				if cardPriceAvailabilityData.Products[i].Price <= cardPriceAvailabilityData.Products[j].Price && cardPriceAvailabilityData.Products[i].Price > 0 && cardPriceAvailabilityData.Products[i].PriceRead != "???" {
					swapCard := cardPriceAvailabilityData.Products[i]
					cardPriceAvailabilityData.Products[i] = cardPriceAvailabilityData.Products[j]
					cardPriceAvailabilityData.Products[j] = swapCard
				}
			}
		}
		return &cardPriceAvailabilityData.Products[0], &cardPriceAvailabilityData.Products[1], &cardPriceAvailabilityData.Products[2]
	}

	return nil, nil, nil
}
