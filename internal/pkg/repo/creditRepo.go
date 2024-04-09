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
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

var creditsRepo *CreditsRepo

type CreditsRepo struct {
	client *dynamodb.Client
	Repository
}

func NewCreditsRepo(client *dynamodb.Client) *CreditsRepo {
	once.Do(func() {
		creditsRepo = &CreditsRepo{client: client}
	})
	return creditsRepo
}

func (c *CreditsRepo) CreateUserCredit(ctx *gin.Context, credit *models.Credits) error {
	item, _ := attributevalue.MarshalMap(credit)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableCredits),
		Item:      item,
	}

	_, er := c.client.PutItem(ctx, params)

	return er
}

func (c *CreditsRepo) FetchCreditByUserID(ctx *gin.Context, userID string) (*models.Credits, error) {
	queryInput := dtos.DbQueryInputConditions{
		Index: models.IndexTableCreditsIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnCreditsUserID: userID,
		},
	}

	creditItems, err := c.QueryItems(ctx, models.TableCredits, queryInput)
	if err != nil {
		log.Printf("error fetching column item: %v\n", err)
		return nil, err
	}

	if len(creditItems) != 1 {
		log.Println("something wrong with credits entries")
		return nil, errors.New("something wrong with credits entries")
	}

	creditsEntity := models.Credits{}
	if err := attributevalue.UnmarshalMap(creditItems[0], &creditsEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return &creditsEntity, nil
}

func (c *CreditsRepo) UpdateCreditsLeftByID(ctx *gin.Context, id string, creditsLeft float64) (*models.Credits, error) {
	queryInput := dtos.DbUpdateQueryConditions{
		Key: map[string]types.AttributeValue{
			models.ColumnCreditsID: &types.AttributeValueMemberS{Value: id},
		},
		ToUpdate: map[string]types.AttributeValue{
			models.ColumnCreditsCreditsLeft: &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%f", creditsLeft),
			},
		},
	}

	// Insert item into the database
	updateItem, err := c.UpdateItem(ctx, models.TableCredits, queryInput)
	if err != nil {
		log.Printf("error inserting item: %v\n", err)
		return nil, err
	}

	creditsEntity := models.Credits{}
	if err := attributevalue.UnmarshalMap(updateItem.Attributes, &creditsEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return &creditsEntity, nil
}
