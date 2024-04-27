package repo

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

var contactsRepo *ContactsRepo

type ContactsRepo struct {
	Repository
}

func newContactsRepos(client *dynamodb.Client) {
	contactsRepo = &ContactsRepo{Repository: Repository{dbClient: client}}
}

func GetContactsRepo() *ContactsRepo {
	return contactsRepo
}

func (c *ContactsRepo) Create(ctx *gin.Context, model *models.Contacts) error {
	item, _ := attributevalue.MarshalMap(model)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableContacts),
		Item:      item,
	}

	_, er := c.dbClient.PutItem(ctx, params)

	return er
}

// BulkCreateContacts inserts multiple contacts in bulk into the DynamoDB table
func (c *ContactsRepo) BulkCreate(ctx *gin.Context, contacts []models.Contacts) error {
	// Prepare the list of PutRequests for batch write operation
	var putRequests []types.WriteRequest
	for _, contact := range contacts {
		item, err := attributevalue.MarshalMap(contact)
		if err != nil {
			return fmt.Errorf("error marshaling contact item: %v", err)
		}

		putRequests = append(putRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		})
	}

	// Split the putRequests into batches of 25 (maximum batch size allowed)
	const batchSize = 25
	numBatches := (len(putRequests) + batchSize - 1) / batchSize

	// Perform batch write operation for each batch
	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(putRequests) {
			end = len(putRequests)
		}

		batchWriteRequests := putRequests[start:end]

		params := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				models.TableContacts: batchWriteRequests,
			},
		}

		_, err := c.dbClient.BatchWriteItem(ctx, params)
		if err != nil {
			return fmt.Errorf("error performing batch write operation: %v", err)
		}
	}

	return nil
}

func (s *ContactsRepo) FetchAllByConditions(ctx *gin.Context, conditions dtos.GetContactsRequest) ([]models.Contacts, error) {
	queryInput := dtos.DbFilterQueryConditions{}
	queryInput.Filters = make(map[string]types.AttributeValue)

	if !utils.IsEmpty(conditions.Name) {
		queryInput.Filters[models.ColumnContactsName] = &types.AttributeValueMemberS{
			Value: conditions.Name,
		}
	}

	if !utils.IsEmpty(conditions.Email) {
		queryInput.Filters[models.ColumnContactsEmail] = &types.AttributeValueMemberS{
			Value: conditions.Email,
		}
	}

	if !utils.IsEmpty(conditions.Mobile) {
		queryInput.Filters[models.ColumnContactsMobile] = &types.AttributeValueMemberS{
			Value: conditions.Mobile,
		}
	}

	if !utils.IsEmpty(conditions.GroupName) {
		queryInput.Filters[models.ColumnContactsGroupName] = &types.AttributeValueMemberS{
			Value: conditions.GroupName,
		}
	}

	if !utils.IsEmpty(conditions.From) {
		queryInput.Filters[fmt.Sprintf("%s_%s", models.ColumnCreatedAt, ComparisonOperatorGE)] = &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d", conditions.From),
		}
	}

	if !utils.IsEmpty(conditions.To) {
		queryInput.Filters[fmt.Sprintf("%s_%s", models.ColumnCreatedAt, ComparisonOperatorLE)] = &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d", conditions.To),
		}
	}

	contactsItems, err := s.ScanItems(ctx, models.TableContacts, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	var contactsEntities []models.Contacts
	if err := attributevalue.UnmarshalListOfMaps(contactsItems, &contactsEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return contactsEntities, nil
}
