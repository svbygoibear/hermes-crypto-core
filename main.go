package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"hermes-crypto-core/internal/db"
	"hermes-crypto-core/internal/handlers/coins"
	"hermes-crypto-core/internal/handlers/users"
	"hermes-crypto-core/internal/middleware"

	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambda

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Default().Println("Could not load .env file, using environment variables.")
	}

	// DB initialization
	db.Init()

	// Set up Gin
	r := gin.Default()
	// Add middleware for panic recovery
	r.Use(middleware.RecoverMiddleware(), middleware.CORSMiddleware())

	// Routes for the users API
	// Votes of users
	r.GET("users/:id/votes", users.GetUserVotesById)
	r.POST("users/:id/votes", users.CreateUserVote)
	r.GET("users/:id/votes/result", users.GetLastUserVoteResult)
	// Health check
	r.GET("users/health", users.HealthCheck)
	// Users base
	r.GET("users", users.GetUsers)
	r.GET("users/:id", users.GetUser)
	r.POST("users", users.CreateUser)
	r.DELETE("users/:id", users.DeleteUser)

	// Routes for the coins API
	// Coin Results
	r.GET("coins/btc", coins.GetCurrentBTCCoinValueInUSD)

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
