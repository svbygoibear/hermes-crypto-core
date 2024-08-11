package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hermes-crypto-core/internal/models"
)

// HealthCheck handles GET requests to check the API health
func HealthCheck(c *gin.Context) {
	statusResponse := models.HealthCheck{
		Status: "healthy",
		Api:    "users",
	}
	c.JSON(http.StatusOK, statusResponse)
}
