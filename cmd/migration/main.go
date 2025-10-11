package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type TableSchema struct {
	TableName              string                       `json:"TableName"`
	KeySchema              []types.KeySchemaElement     `json:"KeySchema"`
	AttributeDefinitions   []types.AttributeDefinition  `json:"AttributeDefinitions"`
	ProvisionedThroughput  *types.ProvisionedThroughput `json:"ProvisionedThroughput"`
	GlobalSecondaryIndexes []types.GlobalSecondaryIndex `json:"GlobalSecondaryIndexes"`
}

func main() {

	logger.Build()

	// Check if the number of command-line arguments is as expected
	if len(os.Args) < 3 {
		fmt.Println("arguments not enough to start the migration")
		os.Exit(1)
	}

	driver := os.Args[1]
	command := os.Args[2]

	if command != "up" {
		fmt.Println("up migrations are allowed only")
		os.Exit(1)
	}

	switch driver {
	case "mysql":
		mysqlMigrations(driver)
		break
	case "dynamo":
		dynamoMigrations()
		break
	}

}

func mysqlMigrations(driver string) {
	// Load config
	config.LoadConfig()

	// Connect to MySQL database
	db, err := sql.Open("mysql", config.GetConfig().Db.URL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a new instance of migrate
	m, err := migrate.New("file://"+utils.GetFilePath("internal/migrations"),
		fmt.Sprintf("%s://", driver)+config.GetConfig().Db.URL())
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Run migrations up to the latest version
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("Migrations completed successfully")
}

func dynamoMigrations() {
	// Load config
	config.LoadConfig()

	// Read table schemas from JSON file
	tableSchemas, err := readTableSchemas(utils.GetFilePath("internal/migrations/tables.json"))
	if err != nil {
		fmt.Println("Error reading table schemas:", err)
		return
	}

	// Create DynamoDB client
	cfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(config.GetConfig().Db.Dynamo.Region),
		awsConfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: config.GetConfig().Db.Dynamo.EndPoint,
				}, nil
			}),
		),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.GetConfig().Db.Dynamo.KeyID,
				config.GetConfig().Db.Dynamo.AccessKey,
				"", // session token, if any
			),
		),
	)

	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return
	}

	client := dynamodb.NewFromConfig(cfg)

	// Create tables
	for _, schema := range tableSchemas {
		input := &dynamodb.CreateTableInput{
			TableName:              &schema.TableName,
			KeySchema:              schema.KeySchema,
			AttributeDefinitions:   schema.AttributeDefinitions,
			ProvisionedThroughput:  schema.ProvisionedThroughput,
			GlobalSecondaryIndexes: schema.GlobalSecondaryIndexes,
		}
		_, err := client.CreateTable(context.Background(), input)
		if err != nil {
			fmt.Printf("Error creating table %s: %v\n", schema.TableName, err)
			continue // skip if this table creation giving and move to next.
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
