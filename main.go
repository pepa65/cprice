package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type MessariResponse struct {
	Data struct {
		Name string `json:"name"`
		Symbol string `json:"symbol"`
		Slug string `json:"slug"`
		MarketData struct {
			Price float64 `json:"price_usd"`
			Pct1H float64 `json:"percent_change_usd_last_1_hour"`
			Pct24H float64 `json:"percent_change_usd_last_24_hours"`
		} `json:"market_data"`
		Marketcap struct {
			Mktcap float64 `json:"current_marketcap_usd"`
		} `json:"marketcap"`
		RoiData struct {
			Pct1W float64 `json:"percent_change_last_1_week"`
			Pct1M float64 `json:"percent_change_last_1_month"`
			Pct3M float64 `json:"percent_change_last_3_months"`
			Pct1Y float64 `json:"percent_change_last_1_year"`
		} `json:"roi_data"`
	} `json:"data"`
}

var self = ""

func main() {
	apikey := os.Getenv("CPRICE_API")
	for _, arg := range os.Args {
		if self == "" { // Get binary name (arg0)
			selves := strings.Split(arg, "/")
			self = selves[len(selves)-1]
			continue
		}
		//res, err := http.Get("https://data.messari.io/api/v1/assets/" + cryptoName + "/metrics?fields=name,symbol,slug,market_data/price_usd")
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://data.messari.io/api/v1/assets/" + arg + "/metrics", nil)
		if apikey != "" {
			req.Header.Add("x-messari-api-key", apikey)
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Could not reach server")
		}
		defer res.Body.Close()
		if res.StatusCode == 404 {
			fmt.Printf("Unsupported cryptocurrency: %s\n", arg)
			continue
		} else if res.StatusCode < 200 || res.StatusCode > 299 {
			fmt.Printf("API returned status code: %d\n", res.StatusCode)
			if res.StatusCode == 429 {
				fmt.Println("Too many API requests to messari.io")
				return
			}
			continue
		}
		var messariRes MessariResponse
		dec := json.NewDecoder(res.Body)
		err = dec.Decode(&messariRes)
		if err != nil {
			fmt.Println("Invalid JSON response")
		} else {
			fmt.Printf("%16s: USD %7.4f\n", messariRes.Data.Name + " " + messariRes.Data.Symbol, messariRes.Data.MarketData.Price)
		}
	}
}
