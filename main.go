package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// status.timestamp data.name data.symbol data.market_data.price_usd data.market_data.percent_change_usd_last_1_hour data.market_data.percent_change_usd_last_24_hours data.marketcap.current_marketcap_usd data.roi_data.percent_change_last_1_week data.roi_data.percent_change_last_1_month data.roi_data.percent_change_last_3_months data.roi_data.percent_change_last_1_year
type MessariResponse struct {
	Data struct {
		Name string `json:"name"`
		Symbol string `json:"symbol"`
		MarketData struct {
			PriceUsd float64 `json:"price_usd"`
		} `json:"market_data"`
	} `json:"data"`
}

func main() {
	cryptoName := os.Args[1]
	//res, err := http.Get("https://data.messari.io/api/v1/assets/" + cryptoName + "/metrics?fields=name,symbol,slug,market_data/price_usd")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://data.messari.io/api/v1/assets/" + cryptoName + "/metrics?fields=name,symbol,slug,market_data/price_usd", nil)
	req.Header.Add("x-messari-api-key", "2a0a30cc-0bac-4b2b-ab74-670f1beff880")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not reach server")
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		fmt.Printf("Unsupported cryptocurrency: %s\n", cryptoName)
		return
	} else if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Printf("API returned status code: %d\n", res.StatusCode)
		return
	}
	var messariRes MessariResponse
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&messariRes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s: USD %0.10f\n", messariRes.Data.Name, messariRes.Data.Symbol, messariRes.Data.MarketData.PriceUsd)
}
