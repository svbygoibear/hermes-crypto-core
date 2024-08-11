package db

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"hermes-crypto-core/internal/models"
)

var client *dynamodb.Client

const tableName = "hermes-crypto-users"

// Init initializes the DynamoDB client
func Init() {
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	region := os.Getenv("AWS_DYNAMODB_REGION")
	keyId := os.Getenv("AWS_DYNAMODB_ACCESS_KEY_ID")
	accessKey := os.Getenv("AWS_DYNAMODB_SECRET_ACCESS_KEY")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if endpoint != "" {
					return aws.Endpoint{URL: endpoint}, nil
				}
				// Fall back to default endpoint resolution
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			})),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(keyId, accessKey, "")),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Add client configuration
	client = dynamodb.NewFromConfig(cfg)

	log.Default().Print("DynamoDB client created")

	// Test connection before proceeding
	hasTable := tableExists()

	log.Default().Print("Do we have a table? ", hasTable)

	if !hasTable {
		// Ensure the table exists
		createTableIfNotExists()
	}
}

func tableExists() bool {
	existingTables, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return false
	}
	log.Println("Successfully connected to DynamoDB")

	// Print out all the table names
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
func GetAllUsers() ([]models.User, error) {
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
func GetUserByID(id string) (*models.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil // Item not found
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a specific user by Email
func GetUserByEmail(id string) (*models.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"Email": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil // Item not found
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user entry in the DynamoDB table
func CreateUser(user models.User) (*models.User, error) {
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
func UpdateUser(id string, user models.User) (*models.User, error) {
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
func DeleteUser(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := client.DeleteItem(context.TODO(), input)
	return err
}
