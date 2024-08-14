package coin

import (
	"context"
	"fmt"

	"os"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"

	"hermes-crypto-core/internal/models"
)

func BinanceGetCurrentExchangeRate() (*float64, error) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_SECRET_KEY")
	client := binance.NewClient(apiKey, apiSecret)

	prices, err := client.NewListPricesService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		fmt.Print("Something went wrong...")
		return nil, models.ReturnError{ErrorMessage: "Failed to retrieve data from Binance API"}
	}

	if len(prices) > 0 {
		priceFloat, err := strconv.ParseFloat(strings.TrimSpace(prices[0].Price), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse price: %w", err)
		}
		fmt.Printf("BINANCE: Current BTC price: $%.2f\n", priceFloat)
		return &priceFloat, nil
	} else {
		fmt.Println("No price data available")
		zero := 0.0
		return &zero, nil
	}
}
