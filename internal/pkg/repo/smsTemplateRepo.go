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

func (s *SmsTemplateRepo) FetchSmsTemplateByUserIDTemplateID(ctx *gin.Context, userID, templateID string) (*models.SmsTemplate, error) {
	queryInput := dtos.DbConditions{
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

	if len(smsTemplateItems) != 1 {
		log.Println("something wrong with sms template entries")
		return nil, errors.New("something wrong with sms template entries")
	}

	smsTemplateEntity := models.SmsTemplate{}
	if err := attributevalue.UnmarshalMap(smsTemplateItems[0], &smsTemplateEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return &smsTemplateEntity, nil
}
