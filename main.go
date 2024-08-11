package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"hermes-crypto-core/internal/handlers/users"
	"hermes-crypto-core/internal/middleware"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// DB initialization
	// db.Init()

	// Set up Gin
	r := gin.Default()
	r.Use(middleware.RecoverMiddleware())

	// Routes for the users API
	r.GET("users/health", users.HealthCheck)
	// r.GET("users/votes", users.GetItems)
	// r.POST("users/votes", users.CreateItem)

	ginLambda = ginadapter.New(r)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.Path == "/users/health" {
		res := &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "OK",
		}

		return *res, nil
	}

	res := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "NOT WHAT I WAS LOOKING FOR",
	}

	return *res, nil

	// log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
	// response, _ := ginLambda.Proxy(request)
	// return response, nil
}

func main() {
	lambda.Start(handler)
}
