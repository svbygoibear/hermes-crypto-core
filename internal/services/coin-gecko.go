package coin_gecko

import (
	"fmt"
	"os"

	"github.com/JulianToledano/goingecko"

	"hermes-crypto-core/internal/models"
)

// CoinGeckoService is a service that interacts with the CoinGecko API
func GetCurrentExchangeRate() (*float64, error) {
	apiKey := os.Getenv("GECKO_API_KEY")
	cgClient := goingecko.NewClient(nil, apiKey)
	defer cgClient.Close()

	data, err := cgClient.CoinsId("bitcoin", true, true, true, false, false, false)
	if err != nil {
		fmt.Print("Something went wrong...")
		return nil, models.ReturnError{ErrorMessage: "Failed to retrieve data from CoinGecko API"}
	}
	fmt.Printf("Bitcoin price is: %f$", data.MarketData.CurrentPrice.Usd)

	return &data.MarketData.CurrentPrice.Usd, nil
}
