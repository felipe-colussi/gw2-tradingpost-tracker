package main

import (
	"fmt"
	"gw2-tradingpost-tracker/runes"

	"os"
)

func main() {
	apiKey, _ := os.LookupEnv("API_KEY")
	r, err := runes.New(apiKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ConvertRunesIntoCharmsAndSell(true, true, true)
}
