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

var s3ProcessingRepo *S3ProcessingRepo

type S3ProcessingRepo struct {
	Repository
}

func newS3ProcessingRepo(client *dynamodb.Client) {
	s3ProcessingRepo = &S3ProcessingRepo{Repository: Repository{dbClient: client}}
}

func GetS3ProcessingRepo() *S3ProcessingRepo {
	return s3ProcessingRepo
}

func (c *S3ProcessingRepo) Create(ctx *gin.Context, model *models.S3Processing) error {
	item, _ := attributevalue.MarshalMap(model)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableS3Processing),
		Item:      item,
	}

	_, er := c.dbClient.PutItem(ctx, params)

	return er
}

func (s *S3ProcessingRepo) FetchByName(ctx *gin.Context, name string) ([]models.S3Processing, error) {
	queryInput := dtos.DbQueryInputConditions{
		Index: models.IndexTableS3ProcessingIndexName,
		PKey: map[string]interface{}{
			models.ColumnS3ProcessingName: name,
		},
	}

	s3ProcessingItems, err := s.QueryItems(ctx, models.TableS3Processing, queryInput)
	if err != nil {
		log.Printf("error fetching s3Processing item: %v\n", err)
		return nil, err
	}

	if len(s3ProcessingItems) == 0 {
		log.Println("no data found")
		return nil, nil
	}

	var s3ProcessingItemsEntity []models.S3Processing
	if err := attributevalue.UnmarshalListOfMaps(s3ProcessingItems, &s3ProcessingItemsEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return s3ProcessingItemsEntity, nil
}

func (c *S3ProcessingRepo) UpdateStatusByID(ctx *gin.Context, id string, status string) (*models.S3Processing, error) {
	queryInput := dtos.DbUpdateQueryConditions{
		Key: map[string]types.AttributeValue{
			models.ColumnS3ProcessingId: &types.AttributeValueMemberS{Value: id},
		},
		ToUpdate: map[string]types.AttributeValue{
			models.ColumnS3ProcessingStatus: &types.AttributeValueMemberS{Value: status},
		},
	}

	// Insert item into the database
	updateItem, err := c.UpdateItem(ctx, models.TableS3Processing, queryInput)
	if err != nil {
		log.Printf("error inserting item: %v\n", err)
		return nil, err
	}

	s3ProcessingEntity := models.S3Processing{}
	if err := attributevalue.UnmarshalMap(updateItem.Attributes, &s3ProcessingEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s3ProcessingEntity, nil
}

func (s *S3ProcessingRepo) FetchAllByConditions(ctx *gin.Context, conditions dtos.GetS3ProcessingRequest) ([]models.S3Processing, error) {
	queryInput := dtos.DbFilterQueryConditions{}
	queryInput.Filters = make(map[string]types.AttributeValue)

	if !utils.IsEmpty(conditions.Name) {
		queryInput.Filters[models.ColumnS3ProcessingName] = &types.AttributeValueMemberS{
			Value: conditions.Name,
		}
	}

	if !utils.IsEmpty(conditions.Status) {
		queryInput.Filters[models.ColumnS3ProcessingStatus] = &types.AttributeValueMemberS{
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

	s3ProcessingItems, err := s.ScanItems(ctx, models.TableS3Processing, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	var s3ProcessingEntities []models.S3Processing
	if err := attributevalue.UnmarshalListOfMaps(s3ProcessingItems, &s3ProcessingEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return s3ProcessingEntities, nil
}
