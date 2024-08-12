package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"hermes-crypto-core/internal/coin"
	"hermes-crypto-core/internal/db"
	"hermes-crypto-core/internal/handlers/users"
	"hermes-crypto-core/internal/middleware"

	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambda

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Default().Println("Error loading .env file")
	}

	// DB initialization
	db.Init()

	currentExchangeRate, err := coin.GetCurrentExchangeRate()
	if err != nil {
		log.Fatalf("Error getting current exchange rate: %v", err)
	}
	log.Printf("Current exchange rate: %f", *currentExchangeRate)

	// Set up Gin
	r := gin.Default()
	// Add middleware for panic recovery
	r.Use(middleware.RecoverMiddleware())

	// Routes for the users API
	r.GET("users/health", users.HealthCheck)
	r.GET("users", users.GetUsers)
	r.GET("users/:id/vote", users.GetUserAndVotes)
	r.POST("users/vote", users.CreateUserAndVotes)

	// Set up the Lambda proxy
	ginLambda = ginadapter.New(r)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, _ := ginLambda.Proxy(request)
	return response, nil
}

func main() {
	lambda.Start(handler)
}
