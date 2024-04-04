package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	once     sync.Once
	baseRepo *Repository
)

type Repository struct {
	dbClient *dynamodb.Client
}

func NewRepository(dbClient *dynamodb.Client) {
	once.Do(func() {
		baseRepo = &Repository{dbClient: dbClient}
	})
}

func GetRepository() *Repository {
	return baseRepo
}

// CreateItem creates a new item in the database
func (r *Repository) CreateItem(ctx context.Context, tableName string, item map[string]types.AttributeValue) error {
	// Define the input for PutItem operation
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: &tableName,
	}

	// Execute PutItem operation
	_, err := r.dbClient.PutItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// GetItemByID retrieves an item by its ID from the database
func (r *Repository) GetItemByID(ctx context.Context, tableName string, key map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	// Implement logic to fetch an item by ID from the database
	return nil, errors.New("not implemented")
}

// UpdateItem updates an existing item in the database
func (r *Repository) UpdateItem(ctx context.Context, tableName string, key map[string]types.AttributeValue, updateExpression string, expressionAttributeValues map[string]types.AttributeValue) error {
	// Implement logic to update an item in the database
	return nil
}

// DeleteItem deletes an item from the database
func (r *Repository) DeleteItem(ctx context.Context, tableName string, key map[string]types.AttributeValue) error {
	// Implement logic to delete an item from the database
	return nil
}
