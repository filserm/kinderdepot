package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type StockData struct {
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"currentPrice"`
}

var stocks = map[string]string{
	"covestro":    "COV.DE",
	"bayer":       "BAYN.DE",
	"bmw":         "BMW.DE",
	"continental": "CON.DE",
	"porsche":     "P911.DE",
	"msci world":  "MSCI",
	"mtu":         "MTU.DE",
	"infineon":    "IFX.DE",
	"linde":       "LIN.DE",
}

func fetchStockPrice(symbol string) (float64, error) {
    apiKey := os.Getenv("API_KEY")
    if apiKey == "" {
        return 0, fmt.Errorf("API_KEY not set")
    }

    apiURL := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, apiKey)
    resp, err := http.Get(apiURL)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return 0, fmt.Errorf("API response: %s", resp.Status)
    }

    var data map[string]map[string]string
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return 0, err
    }

    globalQuote, ok := data["Global Quote"]
    if !ok {
        return 0, fmt.Errorf("no Global Quote found in response")
    }

    priceStr, ok := globalQuote["05. price"]
    if !ok {
        return 0, fmt.Errorf("no price found in Global Quote")
    }

    price, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
        return 0, fmt.Errorf("error converting price to float: %v", err)
    }

    return price, nil
}


func stocksHandler(w http.ResponseWriter, r *http.Request) {
	var results []StockData
	for name, symbol := range stocks {
		price, err := fetchStockPrice(symbol)
		if err != nil {
			log.Println("Error fetching price for", name, ":", err)
			continue
		}
		results = append(results, StockData{
			Symbol:       name,
			CurrentPrice: price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, falling back to system environment variables")
	}

	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoint
	http.HandleFunc("/api/stocks", stocksHandler)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
