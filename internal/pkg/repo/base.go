package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
)

const (
	ComparisonOperatorEQ = "eq"
	ComparisonOperatorL  = "l"
	ComparisonOperatorG  = "g"
	ComparisonOperatorLE = "le"
	ComparisonOperatorGE = "ge"
)

var (
	once     sync.Once
	baseRepo *Repository

	listOfComparisonOperator = map[string]string{
		ComparisonOperatorEQ: "=",
		ComparisonOperatorL:  "<",
		ComparisonOperatorG:  ">",
		ComparisonOperatorLE: "<=",
		ComparisonOperatorGE: ">=",
	}
)

func InitialiseRepositories(dynamodb *dynamodb.Client) {
	once.Do(func() {
		newRepository(dynamodb)
		newCreditsAuditRepo(dynamodb)
		newCreditsRepo(dynamodb)
		newGroupRepo(dynamodb)
		newPendingJobsRepo(dynamodb)
		newS3ProcessingRepo(dynamodb)
		newContactsRepos(dynamodb)
		newOtpRepo(dynamodb)
		newPricingRepo(dynamodb)
		newPromotionRepo(dynamodb)
		newSmsAuditRepo(dynamodb)
		newSmsCampaignRepo(dynamodb)
		newSmsSenderRepo(dynamodb)
		newSmsTemplateRepo(dynamodb)
		newSubscriptionRepo(dynamodb)
		newSubscriptionRepo(dynamodb)
		newUserRepo(dynamodb)
		newPaymentsRepo(dynamodb)
	})
}

type Repository struct {
	dbClient *dynamodb.Client
}

func newRepository(dbClient *dynamodb.Client) {
	baseRepo = &Repository{dbClient: dbClient}
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
func (r *Repository) UpdateItem(ctx context.Context, tableName string, dbUpdateQueryConditions dtos.DbUpdateQueryConditions) (*dynamodb.UpdateItemOutput, error) {
	// Extract key and attributes to update from inputMap
	keyAttributes := dbUpdateQueryConditions.Key

	updateAttributesMap := dbUpdateQueryConditions.ToUpdate

	// Create an UpdateItemInput object with the provided parameters
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       keyAttributes,
		UpdateExpression:          aws.String("SET " + generateUpdateExpressionForUpdate(updateAttributesMap)),
		ExpressionAttributeNames:  generateUpdateExpressionNameForUpdate(updateAttributesMap),
		ExpressionAttributeValues: generateExpressionAttributeValuesForUpdate(updateAttributesMap),
	}

	// Send the UpdateItem request
	updateItem, err := r.dbClient.UpdateItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update item in DynamoDB: %w", err)
	}

	return updateItem, nil
}

// DeleteItem deletes an item from the database
func (r *Repository) DeleteItem(ctx context.Context, tableName string, key map[string]types.AttributeValue) error {
	// Implement logic to delete an item from the database
	return nil
}

// QueryItems retrieves items from the DynamoDB table based on the provided query condition.
func (r *Repository) QueryItems(ctx context.Context, tableName string, conditions dtos.DbQueryInputConditions) ([]map[string]types.AttributeValue, error) {
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

	if !utils.IsEmpty(conditions.Index) {
		input.IndexName = aws.String(conditions.Index)
	}

	// Set the limit
	if !utils.IsEmpty(conditions.Limit) {
		input.Limit = aws.Int32(int32(conditions.Limit))
	}

	input.ExpressionAttributeNames = expressionAttributeNames
	input.ExpressionAttributeValues = expressionAttributeValues

	resp, err := r.dbClient.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error querying DynamoDB: %w", err)
	}

	return resp.Items, nil
}

// ScanItems retrieves items from the DynamoDB table based on the provided scan conditions.
func (r *Repository) ScanItems(ctx context.Context, tableName string, conditions dtos.DbFilterQueryConditions) ([]map[string]types.AttributeValue, error) {
	// Prepare ScanInput
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Set FilterExpression if present
	if len(conditions.Filters) > 0 {
		filterExpr, expressionAttributeNames, expressionAttributeValues, err := buildScanFilterExpression(conditions.Filters)
		if err != nil {
			return nil, fmt.Errorf("error building filter expression: %w", err)
		}
		input.FilterExpression = aws.String(filterExpr)
		input.ExpressionAttributeNames = expressionAttributeNames
		input.ExpressionAttributeValues = expressionAttributeValues
	}

	// Execute the Scan operation
	resp, err := r.dbClient.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error scanning DynamoDB: %w", err)
	}

	return resp.Items, nil
}

// buildScanFilterExpression builds the FilterExpression and related attributes for the ScanInput.
func buildScanFilterExpression(filters map[string]types.AttributeValue) (string, map[string]string, map[string]types.AttributeValue, error) {
	var filterExpr string
	expressionAttributeNames := make(map[string]string)
	expressionAttributeValues := make(map[string]types.AttributeValue)

	// Build FilterExpression and related attributes
	// Example: attribute1 = :value1 AND attribute2 > :value2
	for key, value := range filters {
		var valueMap interface{}
		err := attributevalue.Unmarshal(value, &valueMap)
		if err != nil {
			return "", nil, nil, err
		}

		if utils.IsEmpty(valueMap) {
			continue
		}

		// Extract conditional operator if present
		conditionalOperator := listOfComparisonOperator[ComparisonOperatorEQ] // default is = operator
		if parts := strings.Split(key, "__"); len(parts) == 2 {
			key = parts[0] // Extract actual attribute name

			if _, ok := listOfComparisonOperator[parts[1]]; !ok {
				return "", nil, nil, errors.New("invalid comparison operator passed")
			}

			conditionalOperator = listOfComparisonOperator[parts[1]]
		}

		// Generate placeholder for attribute value
		valuePlaceholder := fmt.Sprintf(":value%d", len(expressionAttributeValues)+1)

		// Append condition to FilterExpression
		if filterExpr != "" {
			filterExpr += " AND "
		}
		filterExpr += fmt.Sprintf("%s %s %s", "#"+key, conditionalOperator, valuePlaceholder)

		// Add attribute name to ExpressionAttributeNames
		expressionAttributeNames["#"+key] = key

		// Add attribute value to ExpressionAttributeValues
		expressionAttributeValues[valuePlaceholder] = value
	}

	return filterExpr, expressionAttributeNames, expressionAttributeValues, nil
}

// buildExpression dynamically builds the KeyConditionExpression, FilterExpression,
// ExpressionAttributeNames, and ExpressionAttributeValues based on the provided conditions.
func buildExpression(conditions dtos.DbQueryInputConditions) (string, string, map[string]string, map[string]types.AttributeValue, error) {
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

// generateUpdateExpression generates the update expression string
func generateUpdateExpressionForUpdate(updateAttributes map[string]types.AttributeValue) string {
	var expr strings.Builder
	for key, value := range updateAttributes {
		if utils.IsEmptyAttributeValue(value) {
			continue
		}

		expr.WriteString("#" + key + " = :" + key + ", ")
	}
	return strings.TrimSuffix(expr.String(), ", ")
}

// generateUpdateExpression generates the update expression string
func generateUpdateExpressionNameForUpdate(updateAttributes map[string]types.AttributeValue) map[string]string {
	var expr = make(map[string]string)
	for key, value := range updateAttributes {
		if utils.IsEmptyAttributeValue(value) {
			continue
		}

		expr["#"+key] = key
	}
	return expr
}

// generateExpressionAttributeValues generates the expression attribute values map
func generateExpressionAttributeValuesForUpdate(updateAttributes map[string]types.AttributeValue) map[string]types.AttributeValue {
	expressionAttributeValues := make(map[string]types.AttributeValue)
	for key, value := range updateAttributes {
		if utils.IsEmptyAttributeValue(value) {
			continue
		}

		expressionAttributeValues[":"+key] = value
	}
	return expressionAttributeValues
}
