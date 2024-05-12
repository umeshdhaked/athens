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

var pendingJobsRepo *PendingJobsRepo

type PendingJobsRepo struct {
	Repository
}

func newPendingJobsRepo(client *dynamodb.Client) {
	pendingJobsRepo = &PendingJobsRepo{Repository: Repository{dbClient: client}}
}

func GetPendingJobsRepo() *PendingJobsRepo {
	return pendingJobsRepo
}

func (c *PendingJobsRepo) Create(ctx *gin.Context, model *models.PendingJobs) error {
	item, _ := attributevalue.MarshalMap(model)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TablePendingJobs),
		Item:      item,
	}

	_, er := c.dbClient.PutItem(ctx, params)

	return er
}

func (s *PendingJobsRepo) FetchByNameAndConditions(ctx *gin.Context, name string, conditions dtos.DbScanQueryConditions) (models.PendingJobs, error) {
	queryInput := dtos.DbQueryInputConditions{
		PKey: map[string]interface{}{
			models.ColumnPendingJobsName: name,
		},
	}

	for k, v := range conditions.Filters {
		queryInput.NonPKey[k] = v
	}

	pendingJobItems, err := s.QueryItems(ctx, models.TablePendingJobs, queryInput)
	if err != nil {
		log.Printf("error fetching pending jobs item: %v\n", err)
		return models.PendingJobs{}, err
	}

	if len(pendingJobItems) != 1 {
		log.Println("fetch by name returned more than 1 items")
		return models.PendingJobs{}, errors.New("fetch by name returned more than 1 items")
	}

	var pendingJobEntity models.PendingJobs
	if err := attributevalue.UnmarshalMap(pendingJobItems[0], &pendingJobEntity); err != nil {
		log.Println(err)
		return models.PendingJobs{}, err
	}

	return pendingJobEntity, nil
}

func (s *PendingJobsRepo) FetchAllByConditions(ctx *gin.Context, conditions dtos.GetPendingJobsRequest) ([]models.PendingJobs, error) {
	queryInput := dtos.DbScanQueryConditions{}
	queryInput.Filters = make(map[string]types.AttributeValue)

	if !utils.IsEmpty(conditions.Name) {
		queryInput.Filters[models.ColumnPendingJobsName] = &types.AttributeValueMemberS{
			Value: conditions.Name,
		}
	}

	if !utils.IsEmpty(conditions.Status) {
		queryInput.Filters[models.ColumnPendingJobsStatus] = &types.AttributeValueMemberS{
			Value: conditions.Status,
		}
	}

	if !utils.IsEmpty(conditions.Type) {
		queryInput.Filters[models.ColumnPendingJobsType] = &types.AttributeValueMemberS{
			Value: conditions.Type,
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

	pendingJobItems, err := s.ScanItems(ctx, models.TablePendingJobs, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	var pendingJobEntities []models.PendingJobs
	if err := attributevalue.UnmarshalListOfMaps(pendingJobItems, &pendingJobEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return pendingJobEntities, nil
}

func (c *PendingJobsRepo) UpdateStatusByName(ctx *gin.Context, name string, status string) (*models.PendingJobs, error) {
	queryInput := dtos.DbUpdateQueryConditions{
		Key: map[string]types.AttributeValue{
			models.ColumnPendingJobsName: &types.AttributeValueMemberS{Value: name},
		},
		ToUpdate: map[string]types.AttributeValue{
			models.ColumnPendingJobsStatus: &types.AttributeValueMemberS{Value: status},
		},
	}

	// Insert item into the database
	updateItem, err := c.UpdateItem(ctx, models.TablePendingJobs, queryInput)
	if err != nil {
		log.Printf("error inserting item: %v\n", err)
		return nil, err
	}

	pendingJobsEntity := models.PendingJobs{}
	if err := attributevalue.UnmarshalMap(updateItem.Attributes, &pendingJobsEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return &pendingJobsEntity, nil
}
