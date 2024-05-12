package repo

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"github.com/gin-gonic/gin"
)

var smsTemplateRepo *SmsTemplateRepo

type SmsTemplateRepo struct {
	Repository
}

func newSmsTemplateRepo(client *dynamodb.Client) {
	smsTemplateRepo = &SmsTemplateRepo{Repository: Repository{dbClient: client}}
}

func GetSmsTemplateRepo() *SmsTemplateRepo {
	return smsTemplateRepo
}

func (s *SmsTemplateRepo) FetchByIDAndUserID(ctx *gin.Context, id, userID string) (*models.SmsTemplate, error) {
	queryInput := dtos.DbQueryInputConditions{
		PKey: map[string]interface{}{
			models.ColumnSmsTemplateID: id,
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsTemplateUserID: userID,
		},
	}

	smsTemplateItems, err := s.QueryItems(ctx, models.TableSmsTemplate, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	if len(smsTemplateItems) == 0 {
		logger.GetLogger().
			Info("no data found")
		return nil, errors.New(ErrCodeNoDataFound)
	}

	var smsTemplateEntities []models.SmsTemplate
	if err := attributevalue.UnmarshalListOfMaps(smsTemplateItems, &smsTemplateEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return &smsTemplateEntities[0], nil
}

func (s *SmsTemplateRepo) FetchByUserIDTemplateCode(ctx *gin.Context, userID, templateCode string) ([]models.SmsTemplate, error) {
	queryInput := dtos.DbQueryInputConditions{
		Index: models.IndexTableSmsTemplateIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsTemplateUserID: userID,
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsTemplateTemplateCode: templateCode,
		},
	}

	smsTemplateItems, err := s.QueryItems(ctx, models.TableSmsTemplate, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	if len(smsTemplateItems) == 0 {
		log.Println("no data found")
		return nil, errors.New(ErrCodeNoDataFound)
	}

	var smsTemplateEntities []models.SmsTemplate
	if err := attributevalue.UnmarshalListOfMaps(smsTemplateItems, &smsTemplateEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return smsTemplateEntities, nil
}
