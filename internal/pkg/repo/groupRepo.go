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

var groupRepo *GroupRepo

type GroupRepo struct {
	Repository
}

func newGroupRepo(client *dynamodb.Client) {
	groupRepo = &GroupRepo{Repository: Repository{dbClient: client}}
}

func GetGroupRepo() *GroupRepo {
	return groupRepo
}

func (c *GroupRepo) CreateGroup(ctx *gin.Context, model *models.Group) error {
	item, _ := attributevalue.MarshalMap(model)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableGroup),
		Item:      item,
	}

	_, er := c.dbClient.PutItem(ctx, params)

	return er
}

func (s *GroupRepo) FetchByUserIDAndName(ctx *gin.Context, userID, name string) ([]models.Group, error) {
	if utils.IsEmpty(userID) {
		return nil, errors.New("userId empty error")
	}

	queryInput := dtos.DbQueryInputConditions{
		Index: models.IndexTableGroupIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnGroupUserID: userID,
		},
		NonPKey: map[string]interface{}{
			models.ColumnGroupName: name,
		},
	}

	groupItems, err := s.QueryItems(ctx, models.TableGroup, queryInput)
	if err != nil {
		log.Printf("error fetching group items: %v\n", err)
		return nil, err
	}

	if len(groupItems) == 0 {
		log.Println("no data found")
		return nil, nil
	}

	var groupEntity []models.Group
	if err := attributevalue.UnmarshalListOfMaps(groupItems, &groupEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return groupEntity, nil
}

func (s *GroupRepo) FetchAllByConditions(ctx *gin.Context, conditions dtos.GetGroupRequest) ([]models.Group, error) {
	queryInput := dtos.DbFilterQueryConditions{}
	queryInput.Filters = make(map[string]types.AttributeValue)

	if !utils.IsEmpty(conditions.Name) {
		queryInput.Filters[models.ColumnGroupName] = &types.AttributeValueMemberS{
			Value: conditions.Name,
		}
	}

	if !utils.IsEmpty(conditions.UserID) {
		queryInput.Filters[models.ColumnGroupUserID] = &types.AttributeValueMemberS{
			Value: conditions.UserID,
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

	groupItems, err := s.ScanItems(ctx, models.TableGroup, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	var groupEntities []models.Group
	if err := attributevalue.UnmarshalListOfMaps(groupItems, &groupEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return groupEntities, nil
}
