package coins

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hermes-crypto-core/internal/coin"
	con "hermes-crypto-core/internal/constants"
	"hermes-crypto-core/internal/models"
)

// GetCurrentBTCCoinValue handles GET requests to retrieve the value for the current Bitcoin coin
func GetCurrentBTCCoinValueInUSD(c *gin.Context) {
	currentExchangeRate, err := coin.GetCurrentExchangeRate()
	if err != nil {
		c.JSON(http.StatusFailedDependency, gin.H{"error": "Could not determine current exchange rate", "message": err.Error()})
		return
	}

	coinResult := models.CoinResult{Coin: con.COIN_TYPE_BTC, CoinValue: *currentExchangeRate, CoinValueCurrency: con.COIN_CURRENCY_USD, QueryTime: models.TimestampTime{Time: time.Now()}}

	c.JSON(http.StatusOK, coinResult)
}
