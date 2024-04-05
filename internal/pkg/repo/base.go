package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/models/dtos"
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

// QueryItems retrieves items from the DynamoDB table based on the provided query condition.
func (r *Repository) QueryItems(ctx context.Context, tableName string, conditions dtos.DbConditions) ([]map[string]types.AttributeValue, error) {
	// Build KeyConditionExpression and FilterExpression dynamically
	keyConditionExpr, filterExpr, expressionAttributeNames, expressionAttributeValues, err := buildExpression(conditions)
	if err != nil {
		return nil, err
	}

	// Execute the Query operation
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
	}

	if !utils.IsEmpty(keyConditionExpr) {
		input.KeyConditionExpression = aws.String(keyConditionExpr)

	}

	if !utils.IsEmpty(filterExpr) {
		input.FilterExpression = aws.String(filterExpr)
	}

	input.ExpressionAttributeNames = expressionAttributeNames
	input.ExpressionAttributeValues = expressionAttributeValues

	input.IndexName = aws.String(conditions.Index)

	resp, err := r.dbClient.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error querying DynamoDB: %w", err)
	}

	return resp.Items, nil
}

// buildExpression dynamically builds the KeyConditionExpression, FilterExpression,
// ExpressionAttributeNames, and ExpressionAttributeValues based on the provided conditions.
func buildExpression(conditions dtos.DbConditions) (string, string, map[string]string, map[string]types.AttributeValue, error) {
	var (
		keyConditionExpr string
		filterExpr       string
		err              error
	)
	expressionAttributeNames := make(map[string]string)
	expressionAttributeValues := make(map[string]types.AttributeValue)

	// Build KeyConditionExpression
	pKeyConditions := conditions.PKey
	if len(pKeyConditions) > 0 {
		keyConditionExpr, expressionAttributeNames, expressionAttributeValues, err = buildKeyCondition(pKeyConditions, expressionAttributeNames, expressionAttributeValues)
		if err != nil {
			return "", "", nil, nil, err
		}
	}

	// Build FilterExpression
	nonPKeyConditions := conditions.NonPKey
	if len(nonPKeyConditions) > 0 {
		filterExpr, expressionAttributeNames, expressionAttributeValues, err = buildFilterExpression(nonPKeyConditions, expressionAttributeNames, expressionAttributeValues)
		if err != nil {
			return "", "", nil, nil, err
		}
	}

	return keyConditionExpr, filterExpr, expressionAttributeNames, expressionAttributeValues, nil
}

// buildKeyCondition builds the KeyConditionExpression based on the provided primary key conditions.
func buildKeyCondition(conditions map[string]interface{}, expressionAttributeNames map[string]string, expressionAttributeValues map[string]types.AttributeValue) (string, map[string]string, map[string]types.AttributeValue, error) {
	var keyConditionExpr []string

	// Iterate through the primary key conditions
	for attrName, attrValue := range conditions {
		alias := "#" + attrName

		// Construct the condition expression for each primary key attribute
		condition := fmt.Sprintf("%s = :%s", alias, attrName)

		// Add the condition to the list of key conditions
		keyConditionExpr = append(keyConditionExpr, condition)

		// Add the attribute name to the expression attribute names map
		expressionAttributeNames[alias] = attrName

		// Add the attribute value to the expression attribute values map
		expressionAttributeValues[":"+attrName] = &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%v", attrValue),
		}
	}

	// Join the key condition expressions with "AND"
	expr := strings.Join(keyConditionExpr, " AND ")

	return expr, expressionAttributeNames, expressionAttributeValues, nil
}

// buildFilterExpression builds the FilterExpression based on the provided non-primary key conditions.
func buildFilterExpression(conditions map[string]interface{}, expressionAttributeNames map[string]string, expressionAttributeValues map[string]types.AttributeValue) (string, map[string]string, map[string]types.AttributeValue, error) {
	var filterExpr []string

	// Iterate through the non-primary key conditions
	for attrName, attrValue := range conditions {
		alias := "#" + attrName

		// Construct the condition expression for each non-primary key attribute
		condition := fmt.Sprintf("%s = :%s", alias, attrName)

		// Add the condition to the list of filter conditions
		filterExpr = append(filterExpr, condition)

		// Add the attribute name to the expression attribute names map
		expressionAttributeNames[alias] = attrName

		// Add the attribute value to the expression attribute values map
		expressionAttributeValues[":"+attrName] = &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%v", attrValue),
		}
	}

	// Join the filter expressions with "AND"
	expr := strings.Join(filterExpr, " AND ")

	return expr, expressionAttributeNames, expressionAttributeValues, nil
}
