package sms

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/http"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"golang.org/x/net/html"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	baseRepo *repo.Repository
}

var (
	once    sync.Once
	service *Service
)

func InitialiseService() {
	once.Do(func() {
		service = &Service{
			baseRepo: repo.GetRepository(),
		}
	})
}

func GetService() *Service {
	return service
}

func (s *Service) AddSenderID(c *gin.Context, request dtos.PostSenderIDRequest) (interface{}, error) {
	//  add validation for same name already existed
	conditions := dtos.DbConditions{
		Index: models.IndexTableSmsSenderIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsSenderUserID: "USERID", // ToDO get userid
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsSenderCode: request.SenderCode,
		},
	}

	items, err := s.baseRepo.QueryItems(context.Background(), models.TableSmsSender, conditions)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	if len(items) > 0 {
		log.Fatalf("duplicate UserID/Name entry")
	}

	entrySmsSenderItem := models.SmsSender{
		ID:       uuid.New().String(),
		Code:     request.SenderCode,
		UserID:   "USERID", //TODO fetch from jwt token
		Type:     request.Type,
		Language: request.Language,
		Status:   models.SmsSenderStateCreated,
	}

	entrySmsSenderItemMap, err := attributevalue.MarshalMap(entrySmsSenderItem)
	if err != nil {
		log.Fatalf("error marhsalling item: %v", err)
	}

	// Insert item into the database
	err = s.baseRepo.CreateItem(context.Background(), models.TableSmsSender, entrySmsSenderItemMap)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

func (s *Service) GetSenderID(c *gin.Context, request dtos.GetSenderIDRequest) (interface{}, error) {
	filters := dtos.DbScanQueryConditions{
		Filters: map[string]types.AttributeValue{
			models.ColumnSmsSenderCode: utils.Ternary(!utils.IsEmpty(request.SenderCode), &types.AttributeValueMemberS{
				Value: request.SenderCode,
			}, &types.AttributeValueMemberS{}).(*types.AttributeValueMemberS),
			models.ColumnSmsSenderType: utils.Ternary(!utils.IsEmpty(request.Type), &types.AttributeValueMemberS{
				Value: request.Type,
			}, &types.AttributeValueMemberS{}).(*types.AttributeValueMemberS),
			models.ColumnSmsSenderUserID: utils.Ternary(!utils.IsEmpty(request.UserID), &types.AttributeValueMemberS{
				Value: request.UserID,
			}, &types.AttributeValueMemberS{}).(*types.AttributeValueMemberS),
			models.ColumnSmsSenderStatus: utils.Ternary(!utils.IsEmpty(request.Status), &types.AttributeValueMemberS{
				Value: request.Status,
			}, &types.AttributeValueMemberS{}).(*types.AttributeValueMemberS),
		},
	}

	items, err := s.baseRepo.ScanItems(context.Background(), models.TableSmsSender, filters)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	return items, nil
}

func (s *Service) ApproveSenderID(c *gin.Context, request dtos.ApproveSenderIDRequest) (interface{}, error) {
	//  Check if sender id exists
	conditions := dtos.DbConditions{
		Index: models.IndexTableSmsSenderIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsSenderUserID: request.UserID,
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsSenderCode: request.SenderCode,
		},
	}

	items, err := s.baseRepo.QueryItems(context.Background(), models.TableSmsSender, conditions)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	if len(items) != 1 {
		log.Fatalf("something wrong with UserID/Name entries")
	}

	senderIDEntity := items[0]

	// Update conditions
	updateConditions := dtos.DbUpdateQueryConditions{
		Key: map[string]types.AttributeValue{
			models.ColumnSmsSenderID: &types.AttributeValueMemberS{Value: senderIDEntity["ID"].(*types.AttributeValueMemberS).Value},
		},
		ToUpdate: map[string]types.AttributeValue{
			models.ColumnSmsSenderStatus: &types.AttributeValueMemberS{Value: models.SmsSenderStateApproved},
		},
	}

	// Insert item into the database
	err = s.baseRepo.UpdateItem(context.Background(), models.TableSmsSender, updateConditions)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

func (s *Service) DeActivateSenderID(c *gin.Context, request dtos.DeleteSenderIDRequest) (interface{}, error) {
	//  Check if sender id exists
	conditions := dtos.DbConditions{
		Index: models.IndexTableSmsSenderIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsSenderUserID: "USERID", // TODO get user id from token
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsSenderCode: request.SenderCode,
		},
	}

	items, err := s.baseRepo.QueryItems(context.Background(), models.TableSmsSender, conditions)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	if len(items) != 1 {
		log.Fatalf("sender id does not exists")
	}

	senderIDEntity := items[0]

	// Update conditions
	updateConditions := dtos.DbUpdateQueryConditions{
		Key: map[string]types.AttributeValue{
			models.ColumnSmsSenderID: &types.AttributeValueMemberS{Value: senderIDEntity["ID"].(*types.AttributeValueMemberS).Value},
		},
		ToUpdate: map[string]types.AttributeValue{
			models.ColumnSmsSenderStatus: &types.AttributeValueMemberS{Value: models.SmsSenderStateDeActivated},
		},
	}

	// Insert item into the database
	err = s.baseRepo.UpdateItem(context.Background(), models.TableSmsSender, updateConditions)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

// This function needs a go through since session token is hard coded and this won't work.
func (s *Service) GetSenderIdInfoFromTrai(c *gin.Context, senderID string) (interface{}, error) {
	var (
		prefix string
		name   string
	)

	// Split sender ID into prefix and name
	parts := strings.Split(senderID, "-")
	if len(parts) == 0 || len(parts) > 2 {
		fmt.Println("Invalid sender ID format")
		return nil, errors.New("invalid sender id format")
	}

	prefix = parts[0]
	name = parts[1]

	response, err := http.NewHTTPClient(config.GetConfig().Api.SmsHeader.BaseUrl).
		Method(config.GetConfig().Api.SmsHeader.Method).
		Path(config.GetConfig().Api.SmsHeader.Path).
		Body(map[string]interface{}{
			"prefix": prefix,
			"name":   name,
		}).
		Headers(map[string]string{
			"Cookie": "JSESSIONID=D96F41E701F13898D526ECB13641DE1F.jvm9", // TODO find a way to get this dynamically
		}).
		SetFormData(true).
		Request(c)

	// Parse HTML
	doc, err := html.Parse(strings.NewReader(string(response)))
	if err != nil {
		log.Fatal(err)
	}

	// Find and extract data from <td> tags
	var traiData []string
	var extractData func(*html.Node)
	extractData = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "td" {
			traiData = append(traiData, n.FirstChild.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractData(c)
		}
	}
	extractData(doc)

	if len(traiData) > 2 && traiData[0] == senderID {
		log.Println("Company Profile: ", traiData[1])
		log.Println("Sender id type: ", traiData[2])
	}

	return nil, nil
}
