package repo

import (
	"errors"
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

var cronProcessingRepo *CronProcessingRepo

type CronProcessingRepo struct {
	Repository
}

func newCronProcessingRepo(client *dynamodb.Client) {
	cronProcessingRepo = &CronProcessingRepo{Repository: Repository{dbClient: client}}
}

func GetCronProcessingRepo() *CronProcessingRepo {
	return cronProcessingRepo
}

func (c *CronProcessingRepo) Create(ctx *gin.Context, model *models.CronProcessing) error {
	item, _ := attributevalue.MarshalMap(model)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableCronProcessing),
		Item:      item,
	}

	_, er := c.dbClient.PutItem(ctx, params)

	return er
}

func (c *CronProcessingRepo) FetchByName(ctx *gin.Context, name string) ([]models.CronProcessing, error) {
	queryInput := dtos.DbQueryInputConditions{
		Index: models.IndexTableCronProcessingIndexName,
		PKey: map[string]interface{}{
			models.ColumnCronProcessingName: name,
		},
	}

	items, err := c.QueryItems(ctx, models.TableCronProcessing, queryInput)
	if err != nil {
		log.Printf("error fetching cronProcessing item: %v\n", err)
		return nil, err
	}

	if len(items) == 0 {
		log.Println("no data found")
		return nil, errors.New(ErrCodeNoDataFound)
	}

	var entities []models.CronProcessing
	if err := attributevalue.UnmarshalListOfMaps(items, &entities); err != nil {
		log.Println(err)
		return nil, err
	}

	return entities, nil
}

func (c *CronProcessingRepo) UpdateStatusByID(ctx *gin.Context, id string, status string) (*models.CronProcessing, error) {
	queryInput := dtos.DbUpdateQueryConditions{
		Key: map[string]types.AttributeValue{
			models.ColumnCronProcessingId: &types.AttributeValueMemberS{Value: id},
		},
		ToUpdate: map[string]types.AttributeValue{
			models.ColumnCronProcessingStatus: &types.AttributeValueMemberS{Value: status},
		},
	}

	// Insert item into the database
	updateItem, err := c.UpdateItem(ctx, models.TableCronProcessing, queryInput)
	if err != nil {
		log.Printf("error inserting item: %v\n", err)
		return nil, err
	}

	entity := models.CronProcessing{}
	if err := attributevalue.UnmarshalMap(updateItem.Attributes, &entity); err != nil {
		log.Println(err)
		return nil, err
	}

	return &entity, nil
}

func (c *CronProcessingRepo) FetchAllByConditions(ctx *gin.Context, conditions dtos.GetCronProcessingRequest) ([]models.CronProcessing, error) {
	queryInput := dtos.DbScanQueryConditions{}
	queryInput.Filters = make(map[string]types.AttributeValue)

	if !utils.IsEmpty(conditions.Name) {
		queryInput.Filters[models.ColumnCronProcessingName] = &types.AttributeValueMemberS{
			Value: conditions.Name,
		}
	}

	if !utils.IsEmpty(conditions.Status) {
		queryInput.Filters[models.ColumnCronProcessingStatus] = &types.AttributeValueMemberS{
			Value: conditions.Status,
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

	items, err := c.ScanItems(ctx, models.TableCronProcessing, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	var entities []models.CronProcessing
	if err := attributevalue.UnmarshalListOfMaps(items, &entities); err != nil {
		log.Println(err)
		return nil, err
	}

	return entities, nil
}
