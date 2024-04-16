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
	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/apiClient"
	"github.com/fastbiztech/hastinapura/internal/pkg/http"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/fastbiztech/hastinapura/pkg/http"
	"golang.org/x/net/html"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	baseRepo        *repo.Repository
	creditsRepo     *repo.CreditsRepo
	smsAuditRepo    *repo.SmsAuditRepo
	smsSenderRepo   *repo.SmsSenderRepo
	smsTemplateRepo *repo.SmsTemplateRepo
	smsCampaignRepo *repo.SmsCampaignRepo
}

var (
	once    sync.Once
	service *Service
)

func InitialiseService() {
	once.Do(func() {
		service = &Service{
			baseRepo:        repo.GetRepository(),
			creditsRepo:     repo.GetCreditsRepo(),
			smsAuditRepo:    repo.GetSmsAuditRepo(),
			smsSenderRepo:   repo.GetSmsSenderRepo(),
			smsTemplateRepo: repo.GetSmsTemplateRepo(),
			smsCampaignRepo: repo.GetSmsCampaignRepo(),
		}
	})
}

func GetService() *Service {
	return service
}

func (s *Service) AddSenderCode(c *gin.Context, request dtos.PostSenderCodeRequest) (interface{}, error) {
	// verify sender code duplication check
	smsSenderEntities, err := s.smsSenderRepo.FetchSmsSenderByUserIDSenderCode(c, c.GetString(constants.JwtTokenUserID), request.SenderCode)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if len(smsSenderEntities) > 0 {
		log.Fatalf("duplicate entry")
	}

	entrySmsSenderItem := models.SmsSender{
		ID:     uuid.New().String(),
		Code:   request.SenderCode,
		UserID: c.GetString(constants.JwtTokenUserID), //TODO fetch from jwt token
		Type:   request.Type,
		Status: models.SmsSenderStateCreated,
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

func (s *Service) GetSenderCode(c *gin.Context, request dtos.GetSenderCodeRequest) (interface{}, error) {
	filters := dtos.DbFilterQueryConditions{
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

func (s *Service) ApproveSenderCode(c *gin.Context, request dtos.ApproveSenderCodeRequest) (interface{}, error) {
	//  Check if sender id exists
	smsSenderExistsConditions := dtos.DbQueryInputConditions{
		Index: models.IndexTableSmsSenderIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsSenderUserID: request.UserID,
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsSenderCode: request.SenderCode,
		},
	}

	items, err := s.baseRepo.QueryItems(context.Background(), models.TableSmsSender, smsSenderExistsConditions)
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
	_, err = s.baseRepo.UpdateItem(context.Background(), models.TableSmsSender, updateConditions)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

func (s *Service) DeActivateSenderCode(c *gin.Context, request dtos.DeleteSenderCodeRequest) (interface{}, error) {
	//  Check if sender id exists
	smsSenderExistsConditions := dtos.DbQueryInputConditions{
		Index: models.IndexTableSmsSenderIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsSenderUserID: c.GetString(constants.JwtTokenUserID), // TODO get user id from token
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsSenderCode: request.SenderCode,
		},
	}

	items, err := s.baseRepo.QueryItems(context.Background(), models.TableSmsSender, smsSenderExistsConditions)
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
	_, err = s.baseRepo.UpdateItem(context.Background(), models.TableSmsSender, updateConditions)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

// This function needs a go through since session token is hard coded and this won't work.
func (s *Service) GetSenderCodeInfoFromTrai(c *gin.Context, senderCode string) (interface{}, error) {
	var (
		prefix string
		name   string
	)

	// Split sender ID into prefix and name
	parts := strings.Split(senderCode, "-")
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

	if len(traiData) > 2 && traiData[0] == senderCode {
		log.Println("Company Profile: ", traiData[1])
		log.Println("Sender id type: ", traiData[2])
	}

	return nil, nil
}

func (s *Service) AddSmsTemplate(c *gin.Context, request dtos.PostSmsTemplateRequest) (interface{}, error) {
	// verify sender code
	smsSenderEntities, err := s.smsSenderRepo.FetchSmsSenderByUserIDSenderCode(c, c.GetString(constants.JwtTokenUserID), request.SenderCode)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if smsSenderEntities == nil {
		log.Fatalf("sender code does not exists")
	}

	// Verify Template ID duplication
	smsTemplateEntities, err := s.smsTemplateRepo.FetchSmsTemplateByUserIDTemplateID(c, c.GetString(constants.JwtTokenUserID), request.TemplateID)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if len(smsTemplateEntities) > 0 {
		log.Fatalf("duplicate entry")
	}

	entrySmsTemplateItem := models.SmsTemplate{
		ID:         uuid.New().String(),
		UserID:     c.GetString(constants.JwtTokenUserID),
		SenderID:   smsSenderEntities[0].ID,
		SenderCode: request.SenderCode,
		Body:       request.Body,
		Status:     models.SmsTemplateStateCreated,
		Type:       request.Type,
		Language:   request.Language,
		TemplateID: request.TemplateID,
	}

	entrySmsTemplateItemMap, err := attributevalue.MarshalMap(entrySmsTemplateItem)
	if err != nil {
		log.Fatalf("error marhsalling item: %v", err)
	}

	// Insert item into the database
	err = s.baseRepo.CreateItem(context.Background(), models.TableSmsTemplate, entrySmsTemplateItemMap)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

func (s *Service) GetSmsTemplate(c *gin.Context, request dtos.GetSmsTemplateRequest) (interface{}, error) {
	filters := dtos.DbFilterQueryConditions{
		Filters: map[string]types.AttributeValue{
			models.ColumnSmsTemplateSenderCode: utils.Ternary(!utils.IsEmpty(request.SenderCode), &types.AttributeValueMemberS{
				Value: request.SenderCode,
			}, &types.AttributeValueMemberS{}).(*types.AttributeValueMemberS),
			models.ColumnSmsTemplateTemplateID: utils.Ternary(!utils.IsEmpty(request.TemplateID), &types.AttributeValueMemberS{
				Value: request.TemplateID,
			}, &types.AttributeValueMemberS{}).(*types.AttributeValueMemberS),
			models.ColumnSmsTemplateStatus: utils.Ternary(!utils.IsEmpty(request.Status), &types.AttributeValueMemberS{
				Value: request.Status,
			}, &types.AttributeValueMemberS{}).(*types.AttributeValueMemberS),
			models.ColumnSmsTemplateUserID: utils.Ternary(!utils.IsEmpty(request.UserID), &types.AttributeValueMemberS{
				Value: request.UserID,
			}, &types.AttributeValueMemberS{}).(*types.AttributeValueMemberS),
		},
	}

	items, err := s.baseRepo.ScanItems(context.Background(), models.TableSmsTemplate, filters)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	return items, nil
}

func (s *Service) ApproveSmsTemplate(c *gin.Context, request dtos.ApproveSmsTemplateRequest) (interface{}, error) {
	//  Check if template exists
	smsTemplateExistsCondition := dtos.DbQueryInputConditions{
		//Index: models.IndexTableSmsTemplateIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsTemplateID: request.TemplateID,
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsSenderUserID: request.UserID,
		},
	}

	items, err := s.baseRepo.QueryItems(context.Background(), models.TableSmsTemplate, smsTemplateExistsCondition)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	if len(items) != 1 {
		log.Fatalf("template id does not exists or exists more than 1")
	}

	// Update conditions
	updateConditions := dtos.DbUpdateQueryConditions{
		Key: map[string]types.AttributeValue{
			models.ColumnSmsTemplateID: &types.AttributeValueMemberS{Value: items[0][models.ColumnSmsTemplateID].(*types.AttributeValueMemberS).Value},
		},
		ToUpdate: map[string]types.AttributeValue{
			models.ColumnSmsTemplateStatus: &types.AttributeValueMemberS{Value: models.SmsSenderStateApproved},
		},
	}

	// Insert item into the database
	_, err = s.baseRepo.UpdateItem(context.Background(), models.TableSmsTemplate, updateConditions)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

func (s *Service) UpdateSmsTemplate(c *gin.Context, request dtos.UpdateSmsTemplateRequest) (interface{}, error) {
	//  Check if template exists
	smsTemplateExistsCondition := dtos.DbQueryInputConditions{
		//Index: models.IndexTableSmsTemplateIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSmsTemplateID: request.TemplateID,
		},
		NonPKey: map[string]interface{}{
			models.ColumnSmsSenderUserID: request.UserID,
		},
	}

	items, err := s.baseRepo.QueryItems(context.Background(), models.TableSmsTemplate, smsTemplateExistsCondition)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	if len(items) != 1 {
		log.Fatalf("template id does not exists or exists more than 1")
	}

	// Update conditions
	updateConditions := dtos.DbUpdateQueryConditions{
		Key: map[string]types.AttributeValue{
			models.ColumnSmsTemplateID: &types.AttributeValueMemberS{Value: items[0][models.ColumnSmsTemplateID].(*types.AttributeValueMemberS).Value},
		},
		ToUpdate: map[string]types.AttributeValue{
			models.ColumnSmsTemplateBody:   &types.AttributeValueMemberS{Value: request.Body},
			models.ColumnSmsTemplateStatus: &types.AttributeValueMemberS{Value: request.Status},
		},
	}

	// Insert item into the database
	_, err = s.baseRepo.UpdateItem(context.Background(), models.TableSmsTemplate, updateConditions)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

func (s *Service) DeActivateSmsTemplate(c *gin.Context, request dtos.DeActivateSmsTemplateRequest) (interface{}, error) {
	// Updating status to DEACTIVATED directly. We can implement delete functionality in future if needed.
	_, err := s.UpdateSmsTemplate(c, dtos.UpdateSmsTemplateRequest{
		TemplateID: request.TemplateID,
		UserID:     request.UserID,
		Status:     models.SmsTemplateStateDeActivated,
	})

	if err != nil {
		log.Fatalf("error deactivating sms template: %v", err)
	}

	return nil, nil
}

func (s *Service) SendInstantSms(c *gin.Context, request dtos.PostSmsRequest) (interface{}, error) {
	// Verify Template ID
	smsTemplateEntities, err := s.smsTemplateRepo.FetchSmsTemplateByUserIDTemplateID(c, c.GetString(constants.JwtTokenUserID), request.TemplateID)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if smsTemplateEntities == nil {
		log.Fatalf("sms template does not exists")
	}

	// verify sender code
	smsSenderEntities, err := s.smsSenderRepo.FetchSmsSenderByUserIDSenderCode(c, c.GetString(constants.JwtTokenUserID), request.SenderCode)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if smsSenderEntities == nil {
		log.Fatalf("sender code does not exists")
	}

	// TODO create contact id

	// Send SMS
	_, err = (&apiClient.InstantSmsApiClient{}).SendInstantSms(c)
	if err != nil {
		return nil, err
	}

	// get credits usage
	creditsUsed := getSmsCreditUsage()

	// fetch credit entry
	creditItem, err := s.creditsRepo.FetchCreditByUserID(c, c.GetString(constants.JwtTokenUserID))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Deduct Credits
	creditItem, err = s.creditsRepo.UpdateCreditsLeftByID(c, creditItem.ID, creditItem.CreditsLeft-creditsUsed)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Create sms audit
	entrySmsAuditItem := models.SmsAudit{
		ID:              uuid.New().String(),
		UserID:          c.GetString(constants.JwtTokenUserID), // ToDO get userid from token
		CreditsConsumed: creditsUsed,
		TemplateID:      request.TemplateID,
		SenderCode:      request.SenderCode,
		//ContactID:       "", // TODO
		Status:        models.SmsAuditStatusDelivered,
		TriggeredMode: models.ModeSmsAuditInstant,
	}

	err = s.smsAuditRepo.CreateSmsAudit(c, &entrySmsAuditItem)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return nil, nil
}

func getSmsCreditUsage() float64 {
	return 0.0 // TODO
}

func (s *Service) CreateSmsCampaign(c *gin.Context, request dtos.CreateSmsCampaignRequest) (interface{}, error) {
	// Verify Template ID
	smsTemplateEntities, err := s.smsTemplateRepo.FetchSmsTemplateByUserIDTemplateID(c, c.GetString(constants.JwtTokenUserID), request.TemplateID)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if smsTemplateEntities == nil {
		log.Fatalf("sms template does not exists")
	}

	// verify sender code
	smsSenderEntities, err := s.smsSenderRepo.FetchSmsSenderByUserIDSenderCode(c, c.GetString(constants.JwtTokenUserID), request.SenderCode)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if smsSenderEntities == nil {
		log.Fatalf("sender code does not exists")
	}

	// create sms campaign
	entrySmsCampaignItem := models.SmsCampaign{
		ID:          uuid.New().String(),
		Name:        request.Name,
		ScheduledAt: request.ScheduledAt,
		Status:      models.SmsCampaignStateCreated,
		UserID:      c.GetString(constants.JwtTokenUserID), // TODO get id from token
		TemplateID:  smsTemplateEntities[0].ID,
		SenderCode:  request.SenderCode,
		//Type:        "",
	}

	err = s.smsCampaignRepo.CreateSmsCampaign(c, &entrySmsCampaignItem)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

func (s *Service) GetSmsCampaigns(c *gin.Context, request dtos.GetSmsCampaignsRequest) (interface{}, error) {
	smsCampaignEntities, err := s.smsCampaignRepo.FetchAllByConditions(c, request)
	if err != nil {
		return nil, err
	}
	return smsCampaignEntities, nil
}

func (s *Service) DeActivateSmsCampaign(c *gin.Context, request dtos.DeActivateSmsCampaignRequest) (interface{}, error) {
	// validate sms campaign
	smsCampaignEntity, err := s.smsCampaignRepo.FetchByID(c, request.ID)
	if err != nil {
		return nil, err
	}

	// Update conditions
	_, err = s.smsCampaignRepo.DeleteByID(c, smsCampaignEntity.ID, true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
