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

var smsSenderRepo *SmsSenderRepo

type SmsSenderRepo struct {
	client *dynamodb.Client
	Repository
}

func NewSmsSenderRepo(client *dynamodb.Client) *SmsSenderRepo {
	once.Do(func() {
		smsSenderRepo = &SmsSenderRepo{client: client}
	})

	return smsSenderRepo
}

func (s *SmsSenderRepo) FetchSmsSenderByUserIDSenderCode(ctx *gin.Context, userID, senderCode string) (*models.SmsSender, error) {
	queryInput := dtos.DbConditions{
		Index: models.IndexTableSmsSenderIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsSenderUserID: userID,
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsSenderCode: senderCode,
		},
	}

	smsSenderItems, err := s.QueryItems(ctx, models.TableSmsSender, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	if len(smsSenderItems) != 1 {
		log.Println("something wrong with sms sender entries")
		return nil, errors.New("something wrong with sms sender entries")
	}

	smsSenderEntity := models.SmsSender{}
	if err := attributevalue.UnmarshalMap(smsSenderItems[0], &smsSenderEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return &smsSenderEntity, nil
}
