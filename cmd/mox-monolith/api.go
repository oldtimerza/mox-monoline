package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func FindCheapestCardAtRetailers(cardId int32) *MerchantCardUnit {
	url := fmt.Sprint("https://moxmonolith.com/card/", cardId, "/products?retailers^[^]=2&retailers^[^]=3&retailers^[^]=4&retailers^[^]=6&retailers^[^]=11&retailers^[^]=13&retailers^[^]=15&retailers^[^]=16&retailers^[^]=18&retailers^[^]=19&retailers^[^]=20&retailers^[^]=21&retailers^[^]=22&retailers^[^]=23&retailers^[^]=24&retailers^[^]=25&retailers^[^]=26&retailers^[^]=32&retailers^[^]=33&retailers^[^]=34&retailers^[^]=35")
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

	if len(cardPriceAvailabilityData.Products) > 0 {
		cheapestCard := cardPriceAvailabilityData.Products[0]
		for i := 0; i < len(cardPriceAvailabilityData.Products); i++ {
			if cardPriceAvailabilityData.Products[i].Price <= cheapestCard.Price && cardPriceAvailabilityData.Products[i].Price > 0 && cardPriceAvailabilityData.Products[i].PriceRead != "???" {
				cheapestCard = cardPriceAvailabilityData.Products[i]
			}
		}
		return &cheapestCard
	}

	return nil
}
