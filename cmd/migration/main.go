package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/utils"
)

type TableSchema struct {
	TableName             string                       `json:"TableName"`
	KeySchema             []types.KeySchemaElement     `json:"KeySchema"`
	AttributeDefinitions  []types.AttributeDefinition  `json:"AttributeDefinitions"`
	ProvisionedThroughput *types.ProvisionedThroughput `json:"ProvisionedThroughput"`
}

func main() {

	// Load config
	config.LoadConfig()

	// Read table schemas from JSON file
	tableSchemas, err := readTableSchemas(utils.GetFilePath("internal/migrations/tables.json"))
	if err != nil {
		fmt.Println("Error reading table schemas:", err)
		return
	}

	// Create DynamoDB client
	cfg, err := awsConfig.LoadDefaultConfig(context.Background(),
		awsConfig.WithRegion(config.GetConfig().Aws.Db.Region),
		awsConfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: config.GetConfig().Aws.Db.EndPoint}, nil
			})))

	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return
	}

	client := dynamodb.NewFromConfig(cfg)

	// Create tables
	for _, schema := range tableSchemas {
		input := &dynamodb.CreateTableInput{
			TableName:             &schema.TableName,
			KeySchema:             schema.KeySchema,
			AttributeDefinitions:  schema.AttributeDefinitions,
			ProvisionedThroughput: schema.ProvisionedThroughput,
		}
		_, err := client.CreateTable(context.Background(), input)
		if err != nil {
			fmt.Printf("Error creating table %s: %v\n", schema.TableName, err)
			return
		}
		fmt.Printf("Table created successfully: %s\n", schema.TableName)
	}
}

func readTableSchemas(file string) ([]TableSchema, error) {
	var schemas []TableSchema
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &schemas)
	if err != nil {
		return nil, err
	}
	return schemas, nil
}
