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
	Repository
}

func newSmsSenderRepo(client *dynamodb.Client) {
	smsSenderRepo = &SmsSenderRepo{Repository: Repository{dbClient: client}}
}

func GetSmsSenderRepo() *SmsSenderRepo {
	return smsSenderRepo
}

func (s *SmsSenderRepo) FetchSmsSenderByUserIDSenderCode(ctx *gin.Context, userID, senderCode string) ([]models.SmsSender, error) {
	queryInput := dtos.DbQueryInputConditions{
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

	if len(smsSenderItems) == 0 {
		log.Println("no data found")
		return nil, errors.New(ErrCodeNoDataFound)
	}

	var smsSenderEntity []models.SmsSender
	if err := attributevalue.UnmarshalListOfMaps(smsSenderItems, &smsSenderEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return smsSenderEntity, nil
}
