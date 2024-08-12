package db

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"hermes-crypto-core/internal/models"
)

type dynamoDB struct {
	client *dynamodb.Client
}

var client *dynamodb.Client

const tableName = "hermes-crypto-users"

// Init initializes the DynamoDB client
func Init() {
	dbRegion := os.Getenv("AWS_DYNAMODB_REGION")
	if dbRegion == "" {
		log.Fatal("AWS_DYNAMODB_REGION is not set")
	}

	log.Printf("Initializing DynamoDB client with region: %s", dbRegion)

	var cfg aws.Config
	var err error

	isLocal := os.Getenv("IS_LOCAL")
	if isLocal == "true" {
		// Local development configuration
		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           "http://localhost:1433", // Default local DynamoDB port
				SigningRegion: dbRegion,
			}, nil
		})

		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(dbRegion),
			config.WithEndpointResolver(customResolver),
		)
	} else {
		// Production configuration
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(dbRegion),
		)
	}

	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}

	client = dynamodb.NewFromConfig(cfg)
	DB = &dynamoDB{client: client}

	log.Println("DynamoDB client created successfully")

	// Test connection before proceeding
	hasTable := tableExists()

	log.Printf("Table %s exists: %v", tableName, hasTable)

	if !hasTable {
		createTableIfNotExists()
	}
}

func tableExists() bool {
	existingTables, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return false
	}
	log.Println("Successfully connected to DynamoDB")

	// Gets the table names from the response and matches
	// the table name with the one we are looking for
	for _, table := range existingTables.TableNames {
		var tablePtr *string = &table
		if *tablePtr == tableName {
			log.Println("Table already exists")
			return true
		}
	}

	log.Println("Table not found")
	return false
}

func createTableIfNotExists() {
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		// If the table already exists, ignore the error
		if _, ok := err.(*types.ResourceInUseException); !ok {
			log.Fatalf("Error creating table: %v", err)
		}
	}
}

// GetAllUsers retrieves all users from the DynamoDB table
func (d *dynamoDB) GetAllUsers() ([]models.User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID retrieves a specific user by Id
func (d *dynamoDB) GetUserByID(id string) (*models.User, error) {
	// This is not an ideal solution - this should be optimized in future
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		// KeyConditionExpression: aws.String("Id = :Id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Id": &types.AttributeValueMemberS{Value: id},
		},
		Limit: aws.Int32(1), // We only need one item
	}

	result, err := client.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil // User not found
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a specific user by Email
func (d *dynamoDB) GetUserByEmail(email string) (*models.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"Email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil // User not found
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user entry in the DynamoDB table
func (d *dynamoDB) CreateUser(user models.User) (*models.User, error) {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = client.PutItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user in the DynamoDB table, using their user Id
func (d *dynamoDB) UpdateUser(id string, user models.User) (*models.User, error) {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = client.PutItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser removes a user from the DynamoDB table
func (d *dynamoDB) DeleteUser(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := client.DeleteItem(context.TODO(), input)
	return err
}
