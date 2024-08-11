package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"hermes-crypto-core/internal/db"
	"hermes-crypto-core/internal/handlers/users"
	"hermes-crypto-core/internal/middleware"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// DB initialization
	db.Init()

	// Set up Gin
	r := gin.Default()
	// Add middleware for panic recovery
	r.Use(middleware.RecoverMiddleware())

	// Routes for the users API
	r.GET("users/health", users.HealthCheck)
	r.GET("users/votes", users.GetUsers)
	r.GET("users/votes/:id", users.GetUser)
	r.POST("users/votes", users.CreateUser)

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
