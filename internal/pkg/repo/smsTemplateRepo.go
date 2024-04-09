package repo

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

var smsTemplateRepo *SmsTemplateRepo

type SmsTemplateRepo struct {
	client *dynamodb.Client
	Repository
}

func NewSmsTemplateRepo(client *dynamodb.Client) *SmsTemplateRepo {
	once.Do(func() {
		smsTemplateRepo = &SmsTemplateRepo{client: client}
	})

	return smsTemplateRepo
}

func (s *SmsTemplateRepo) FetchSmsTemplateByUserIDTemplateID(ctx *gin.Context, userID, templateID string) ([]models.SmsTemplate, error) {
	queryInput := dtos.DbQueryInputConditions{
		Index: models.IndexTableSmsTemplateIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsTemplateUserID: userID,
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsTemplateTemplateID: templateID,
		},
	}

	smsTemplateItems, err := s.QueryItems(ctx, models.TableSmsTemplate, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	if len(smsTemplateItems) == 0 {
		log.Println("no data found")
		return nil, errors.New("no data found")
	}

	var smsTemplateEntities []models.SmsTemplate
	if err := attributevalue.UnmarshalListOfMaps(smsTemplateItems, &smsTemplateEntities); err != nil {
		log.Println(err)
		return nil, err
	}

	return smsTemplateEntities, nil
}
