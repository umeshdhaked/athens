package repo

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

var smsCampaignRepo *SmsCampaignRepo

type SmsCampaignRepo struct {
	client *dynamodb.Client
	Repository
}

func NewSmsCampaignRepo(client *dynamodb.Client) *SmsCampaignRepo {
	once.Do(func() {
		smsCampaignRepo = &SmsCampaignRepo{client: client}
	})

	return smsCampaignRepo
}

func (s *SmsCampaignRepo) CreateSmsCampaign(ctx *gin.Context, smsCampaign *models.SmsCampaign) error {
	item, err := attributevalue.MarshalMap(smsCampaign)
	if err != nil {
		log.Printf("error marhsalling item: %v", err)
		return nil
	}

	// Insert item into the database
	err = s.CreateItem(ctx, models.TableSmsCampaign, item)
	if err != nil {
		log.Printf("error inserting item: %v", err)
		return nil
	}

	return err
}

func (s *SmsCampaignRepo) FetchByID(ctx *gin.Context, id string) (models.SmsCampaign, error) {
	queryInput := dtos.DbQueryInputConditions{
		PKey: map[string]interface{}{
			models.ColumnSmsCampaignID: id,
		},
	}

	smsCampaignItems, err := s.QueryItems(ctx, models.TableSmsCampaign, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return models.SmsCampaign{}, err
	}

	if len(smsCampaignItems) != 1 {
		log.Println("fetch by id returned more than 1 items")
		return models.SmsCampaign{}, errors.New("fetch by id returned more than 1 items")
	}

	var smsCampaignEntity models.SmsCampaign
	if err := attributevalue.UnmarshalMap(smsCampaignItems[0], &smsCampaignEntity); err != nil {
		log.Println(err)
		return models.SmsCampaign{}, err
	}

	return smsCampaignEntity, nil
}

func (s *SmsCampaignRepo) FetchSmsCampaignByUserIDAndConditions(ctx *gin.Context, userID string, conditions dtos.DbFilterQueryConditions) ([]models.SmsCampaign, error) {
	queryInput := dtos.DbQueryInputConditions{
		Index: models.IndexTableSmsCampaignIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsTemplateUserID: userID,
		},
	}

	for k, v := range conditions.Filters {
		queryInput.NonPKey[k] = v
	}

	smsCampaignItems, err := s.QueryItems(ctx, models.TableSmsCampaign, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	var smsCampaignEntities []models.SmsCampaign
	if err := attributevalue.UnmarshalListOfMaps(smsCampaignItems, &smsCampaignEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return smsCampaignEntities, nil
}

func (s *SmsCampaignRepo) FetchSmsCampaignByStatusAndConditions(ctx *gin.Context, status string, conditions dtos.DbFilterQueryConditions) ([]models.SmsCampaign, error) {
	queryInput := dtos.DbQueryInputConditions{
		Index: models.IndexTableSmsCampaignIndexStatus,
		PKey: map[string]interface{}{
			models.ColumnSmsTemplateStatus: status,
		},
	}

	for k, v := range conditions.Filters {
		queryInput.NonPKey[k] = v
	}

	smsCampaignItems, err := s.QueryItems(ctx, models.TableSmsCampaign, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	var smsCampaignEntities []models.SmsCampaign
	if err := attributevalue.UnmarshalListOfMaps(smsCampaignItems, &smsCampaignEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return smsCampaignEntities, nil
}

func (s *SmsCampaignRepo) FetchAllByConditions(ctx *gin.Context, conditions dtos.GetSmsCampaignsRequest) ([]models.SmsCampaign, error) {
	queryInput := dtos.DbFilterQueryConditions{}
	queryInput.Filters = make(map[string]types.AttributeValue)

	if !utils.IsEmpty(conditions.Name) {
		queryInput.Filters[models.ColumnSmsCampaignName] = &types.AttributeValueMemberS{
			Value: conditions.Name,
		}
	}

	if !utils.IsEmpty(conditions.UserID) {
		queryInput.Filters[models.ColumnSmsCampaignUserID] = &types.AttributeValueMemberS{
			Value: conditions.UserID,
		}
	}

	if !utils.IsEmpty(conditions.Status) {
		queryInput.Filters[models.ColumnSmsCampaignStatus] = &types.AttributeValueMemberS{
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

	smsCampaignItems, err := s.ScanItems(ctx, models.TableSmsCampaign, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	var smsCampaignEntities []models.SmsCampaign
	if err := attributevalue.UnmarshalListOfMaps(smsCampaignItems, &smsCampaignEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return smsCampaignEntities, nil
}

// TODO soft/hard delete handling
func (c *SmsCampaignRepo) DeleteByID(ctx *gin.Context, id string, soft bool) (*models.SmsCampaign, error) {
	queryInput := dtos.DbUpdateQueryConditions{
		Key: map[string]types.AttributeValue{
			models.ColumnSmsCampaignID: &types.AttributeValueMemberS{Value: id},
		},
		ToUpdate: map[string]types.AttributeValue{
			models.ColumnSmsCampaignStatus: &types.AttributeValueMemberS{
				Value: models.SmsCampaignStateDeActivated,
			},
		},
	}

	// Insert item into the database
	updateItem, err := c.UpdateItem(ctx, models.TableSmsCampaign, queryInput)
	if err != nil {
		log.Printf("error updating item: %v\n", err)
		return nil, err
	}

	smsCampaignEntity := models.SmsCampaign{}
	if err := attributevalue.UnmarshalMap(updateItem.Attributes, &smsCampaignEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return &smsCampaignEntity, nil
}
