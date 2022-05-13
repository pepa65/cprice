package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type MessariJson struct {
	D struct {
		Name string `json:"name"`
		Symbol string `json:"symbol"`
		Slug string `json:"slug"`
		MktData struct {
			Price float64 `json:"price_usd"`
			Pct1H float64 `json:"percent_change_usd_last_1_hour"`
			Pct24H float64 `json:"percent_change_usd_last_24_hours"`
		} `json:"market_data"`
		Mktcap struct {
			Marketcap float64 `json:"current_marketcap_usd"`
		} `json:"marketcap"`
		RoiData struct {
			Pct1W float64 `json:"percent_change_last_1_week"`
			Pct1M float64 `json:"percent_change_last_1_month"`
			Pct3M float64 `json:"percent_change_last_3_months"`
			Pct1Y float64 `json:"percent_change_last_1_year"`
		} `json:"roi_data"`
	} `json:"data"`
}

const (
	version = "0.1.3"
	def     = "\033[0m"
	cya     = "\033[1m\033[36m"
	red     = "\033[1m\033[31m"
	gre     = "\033[1m\033[32m"
)
var (
	self = ""
	apikey = os.Getenv("CPRICE_API")
)

func usage() {
	fmt.Printf("%s%s v%s%s - CLI to access cryptocoin data online\n", gre, self, version, def)
	fmt.Printf("Usage:  %s%s <cryptcurrency>...%s\n", cya, self, def)
	fmt.Printf("messari.io API key: %sCPRICE_API%s=\"%s%s%s\"\n", cya, def, gre, apikey, def)
}

func main() {
	for _, arg := range os.Args {
		if self == "" { // Get binary name (arg0)
			selves := strings.Split(arg, "/")
			self = selves[len(selves)-1]
			if len(os.Args) == 1 {
				usage()
				return
			}
			fmt.Printf("messari.io API key: %sCPRICE_API%s=\"%s%s%s\"\n", cya, def, gre, apikey, def)
			fmt.Println(cya + "            Name Symbol    Value USD    1h  1d  1w  1m  3m  1y  Marketcap" + def)
			continue
		}
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://data.messari.io/api/v1/assets/" + arg + "/metrics", nil)
		if apikey != "" {
			req.Header.Add("x-messari-api-key", apikey)
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(red + "Could not reach server" + def)
			return
		}
		defer res.Body.Close()
		if res.StatusCode == 404 {
			fmt.Printf("%sUnsupported cryptocurrency: %s%s\n", red, def, arg)
			continue
		} else if res.StatusCode < 200 || res.StatusCode > 299 {
			fmt.Printf("%sAPI returned status code: %s%d\n", red, def, res.StatusCode)
			if res.StatusCode == 429 {
				fmt.Println(red + "Too many API requests to messari.io" + def)
				return
			}
			continue
		}
		var m MessariJson
		dec := json.NewDecoder(res.Body)
		err = dec.Decode(&m)
		if err != nil {
			fmt.Println(red + "Invalid JSON response" +def)
			continue
		}
		u := m.D.MktData.Price
		if u == 0 {
			fmt.Printf("%sNo value recorded for %s%s\n", red, def, arg)
			continue
		}
		pct := [6]float64{
			m.D.MktData.Pct1H,
			m.D.MktData.Pct24H,
			m.D.RoiData.Pct1W,
			m.D.RoiData.Pct1M,
			m.D.RoiData.Pct3M,
			m.D.RoiData.Pct1Y,
		}
		change := ""
		for i := 0; i < 6; i++ {
			p := pct[i]
			if p < 0 {
				change += red + fmt.Sprintf(" %3.0f", -p)
			} else {
				change += gre + fmt.Sprintf(" %3.0f", p)
			}
		}
		fmt.Printf("%23s:  %11.5f %s  %.4e\n", m.D.Name + " " + m.D.Symbol, u, change + def, m.D.Mktcap.Marketcap)
	}
}
