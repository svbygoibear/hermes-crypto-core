package coin

import (
	"hermes-crypto-core/internal/models"
)

func GetCurrentExchangeRate() (*float64, error) {
	currentExchangeRate, err := BinanceGetCurrentExchangeRate()
	if err != nil {
		geckoExchangeRate, err := GeckoGetCurrentExchangeRate()
		if err != nil {
			return nil, models.ReturnError{ErrorMessage: "Could not determine current exchange rate"}
		}

		return geckoExchangeRate, nil
	}

	return currentExchangeRate, nil
}
