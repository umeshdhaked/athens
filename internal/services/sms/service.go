package sms

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/apiClient"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/fastbiztech/hastinapura/pkg/http"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"golang.org/x/net/html"
	gormLogger "gorm.io/gorm/logger"

	"github.com/gin-gonic/gin"
)

type Service struct {
	baseRepo        repo.IRepository
	creditsRepo     repo.ICreditsRepo
	smsAuditRepo    repo.ISmsAuditRepo
	smsSenderRepo   repo.ISmsSenderRepo
	smsTemplateRepo repo.ISmsTemplateRepo
	smsCampaignRepo repo.ISmsCampaignRepo
}

var (
	once    sync.Once
	service *Service
)

const (
	CreditUseBelow160Length = 1
	CreditUseAbove160Length = 2
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
	var smsSender models.SmsSender
	err := s.baseRepo.Find(c, &smsSender, map[string]interface{}{
		constants.UserID: c.GetInt64(constants.JwtTokenUserID),
		constants.Code:   request.SenderCode,
	})
	if err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound) {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	if smsSender.ID != 0 {
		logger.GetLogger().Error("sms sender code already exists")
		return nil, err
	}

	entrySmsSenderItem := models.SmsSender{
		Code:   request.SenderCode,
		UserID: c.GetString(constants.JwtTokenUserID),
		Type:   request.Type,
		Status: models.SmsSenderStateCreated,
	}

	// Insert item into the database
	err = s.baseRepo.Create(c, &entrySmsSenderItem)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return nil, nil
}

func (s *Service) GetSenderCode(c *gin.Context, request dtos.GetSenderCodeRequest) (interface{}, error) {
	var smsSenderEntities []models.SmsSender
	err := s.baseRepo.FindMultiplePagination(c, &smsSenderEntities, map[string]interface{}{
		constants.Code:   request.SenderCode,
		constants.Type:   request.Type,
		constants.UserID: request.UserID,
		constants.Status: request.Status,
	}, dtos.Pagination{})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return smsSenderEntities, nil
}

func (s *Service) ApproveSenderCode(c *gin.Context, request dtos.ApproveSenderCodeRequest) (interface{}, error) {
	var smsSender models.SmsSender
	err := s.baseRepo.Find(c, &smsSender, map[string]interface{}{
		constants.UserID: request.UserID,
		constants.Code:   request.SenderCode,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Insert item into the database
	smsSender.Status = models.SmsSenderStateApproved
	err = s.baseRepo.Update(c, &smsSender)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return nil, nil
}

func (s *Service) DeActivateSenderCode(c *gin.Context, request dtos.DeleteSenderCodeRequest) (interface{}, error) {
	//  Check if sender id exists
	var smsSender models.SmsSender
	err := s.baseRepo.Find(c, &smsSender, map[string]interface{}{
		constants.Code: request.SenderCode,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Insert item into the database
	smsSender.Status = models.SmsSenderStateDeActivated
	err = s.baseRepo.Update(c, &smsSender)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
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
	var smsSender models.SmsSender
	err := s.baseRepo.Find(c, &smsSender, map[string]interface{}{
		constants.UserID: c.GetInt64(constants.JwtTokenUserID),
		constants.Code:   request.SenderCode,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Verify Template ID duplication
	var smsTemplate models.SmsTemplate
	err = s.baseRepo.FindMultiple(c, &smsTemplate, map[string]interface{}{
		constants.TemplateCode: request.TemplateCode,
		constants.UserID:       c.GetInt64(constants.JwtTokenUserID),
	})
	if err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound) {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	if smsTemplate.ID != 0 {
		return nil, errors.New("sms template already exists")
	}

	insertSmsTemplate := models.SmsTemplate{
		UserID:       strconv.Itoa(int(c.GetInt64(constants.JwtTokenUserID))),
		SenderID:     smsSender.ID,
		SenderCode:   request.SenderCode,
		Body:         request.Body,
		Status:       models.SmsTemplateStateCreated,
		Type:         request.Type,
		Length:       len(request.Body),
		Language:     request.Language,
		TemplateCode: request.TemplateCode,
	}

	// Insert item into the database
	err = s.baseRepo.Create(c, &insertSmsTemplate)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return nil, nil
}

func (s *Service) GetSmsTemplate(c *gin.Context, request dtos.GetSmsTemplateRequest) (interface{}, error) {

	var smsTemplateEntities []models.SmsTemplate
	err := s.baseRepo.FindMultiplePagination(c, &smsTemplateEntities, map[string]interface{}{
		constants.SenderCode: request.SenderCode,
		constants.ID:         request.TemplateID,
		constants.Status:     request.Status,
		constants.UserID:     request.UserID,
	}, dtos.Pagination{})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return smsTemplateEntities, nil
}

func (s *Service) ApproveSmsTemplate(c *gin.Context, request dtos.ApproveSmsTemplateRequest) (interface{}, error) {
	//  Check if template exists
	var smsTemplate models.SmsTemplate
	err := s.baseRepo.Find(c, &smsTemplate, map[string]interface{}{
		constants.ID:     request.TemplateID,
		constants.UserID: c.GetInt64(constants.JwtTokenUserID),
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Update conditions
	smsTemplate.Status = models.SmsSenderStateApproved
	err = s.baseRepo.Update(c, &smsTemplate)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return smsTemplate, nil
}

func (s *Service) UpdateSmsTemplate(c *gin.Context, request dtos.UpdateSmsTemplateRequest) (interface{}, error) {
	//  Check if template exists
	var smsTemplate models.SmsTemplate
	err := s.baseRepo.Find(c, &smsTemplate, map[string]interface{}{
		constants.ID:     request.TemplateID,
		constants.UserID: request.UserID,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Insert item into the database
	smsTemplate.Body = request.Body
	smsTemplate.Status = request.Status
	err = s.baseRepo.Update(c, &smsTemplate)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return smsTemplate, nil
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
	var smsTemplate models.SmsTemplate
	err := s.baseRepo.Find(c, &smsTemplate, map[string]interface{}{
		constants.ID:     request.TemplateID,
		constants.UserID: c.GetInt64(constants.JwtTokenUserID),
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	if smsTemplate.Status != models.SmsTemplateStateApproved {
		logger.GetLogger().Error("sms template not approved")
		return nil, err
	}

	// verify sender code
	var smsSender models.SmsSender
	err = s.baseRepo.Find(c, &smsSender, map[string]interface{}{
		constants.UserID: c.GetInt64(constants.JwtTokenUserID),
		constants.Code:   request.SenderCode,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
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
	var credit models.Credits
	err = s.baseRepo.Find(c, &credit, map[string]interface{}{
		constants.UserID: c.GetInt64(constants.JwtTokenUserID),
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Deduct Balance
	credit.BalanceLeft = credit.BalanceLeft - creditsUsed
	err = s.baseRepo.Update(c, &credit)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Create sms audit
	smsAuditItem := models.SmsAudit{
		UserID:          strconv.Itoa(int(c.GetInt64(constants.JwtTokenUserID))), // ToDO get userid from token
		CreditsConsumed: creditsUsed,
		TemplateID:      request.TemplateID,
		SenderCode:      request.SenderCode,
		//ContactID:       "", // TODO
		Status:        models.SmsAuditStatusDelivered,
		TriggeredMode: models.ModeSmsAuditInstant,
	}

	err = s.baseRepo.Create(c, &smsAuditItem)
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
	var smsTemplate models.SmsTemplate
	err := s.baseRepo.Find(c, &smsTemplate, map[string]interface{}{
		constants.ID:     request.TemplateID,
		constants.UserID: c.GetInt64(constants.JwtTokenUserID),
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// verify sender code
	var smsSender models.SmsSender
	err = s.baseRepo.Find(c, &smsSender, map[string]interface{}{
		constants.UserID: c.GetInt64(constants.JwtTokenUserID),
		constants.Code:   request.SenderCode,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// create sms campaign
	smsCampaign := models.SmsCampaign{
		Name:        request.Name,
		ScheduledAt: request.ScheduledAt,
		Status:      models.SmsCampaignStateCreated,
		UserID:      strconv.Itoa(int(c.GetInt64(constants.JwtTokenUserID))),
		TemplateID:  smsTemplate.ID,
		SenderCode:  request.SenderCode,
		GroupName:   request.GroupName,
	}

	err = s.baseRepo.Create(c, &smsCampaign)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return nil, nil
}

func (s *Service) GetSmsCampaigns(c *gin.Context, request dtos.GetSmsCampaignsRequest) (interface{}, error) {
	var smsCampaignEntities []models.SmsCampaign
	err := s.baseRepo.FindMultiplePagination(c, &smsCampaignEntities, map[string]interface{}{
		constants.Name:   request.Name,
		constants.Status: request.Status,
		constants.UserID: request.UserID,
	}, dtos.Pagination{
		From: request.From,
		To:   request.To,
	})
	if err != nil {
		return nil, err
	}
	return smsCampaignEntities, nil
}

func (s *Service) DeActivateSmsCampaign(c *gin.Context, request dtos.DeActivateSmsCampaignRequest) (interface{}, error) {
	// validate sms campaign
	var smsCampaign models.SmsCampaign
	err := s.baseRepo.FindByID(c, request.ID, &smsCampaign)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Update conditions
	smsCampaign.Status = models.SmsCampaignStateDeActivated
	err = s.baseRepo.Update(c, &smsCampaign)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return nil, nil
}

// getCampaignsToRun fetches campaigns for which scheduled time has passed
func getCampaignsToRun(ctx *gin.Context) ([]models.SmsCampaign, error) {

	var (
		result []models.SmsCampaign
		err    error
	)

	var smsCampaignEntities []models.SmsCampaign
	err = repo.GetRepository().FindMultiple(ctx, &smsCampaignEntities, map[string]interface{}{
		constants.Status: models.SmsCampaignStateCreated,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	baseTime := time.Now().Unix()
	for _, smsCampaign := range smsCampaignEntities {

		if int(baseTime) > smsCampaign.ScheduledAt {
			result = append(result, smsCampaign)
		}
	}

	return result, nil
}

func isEnoughBalanceForSmsCampaignRun(ctx *gin.Context,
	userID string,
	templateID int64,
	subCategory string,
	batch int) (*models.Credits, float64, bool) {

	perSmsCost, err := getPerSmsCost(ctx, userID, templateID, subCategory)
	if err != nil {
		return nil, 0, false
	}

	var credit models.Credits
	err = repo.GetRepository().Find(ctx, &credit, map[string]interface{}{
		constants.UserID: userID,
	})
	if err != nil {
		logger.GetLogger().
			Info("failed to fetch credit items")
		return nil, 0, false
	}

	if (perSmsCost * float64(batch)) <= credit.Balance {
		return &credit, perSmsCost * float64(batch), true
	}

	return nil, 0, false
}

func getPerSmsCost(ctx *gin.Context,
	userID string,
	templateID int64,
	subCategory string) (float64, error) {

	// fetch sms template
	var smsTemplate models.SmsTemplate
	err := repo.GetRepository().Find(ctx, &smsTemplate, map[string]interface{}{
		constants.ID:     templateID,
		constants.UserID: userID,
	})
	if err != nil {
		logger.GetLogger().
			WithField("error", err).
			Info("failed to fetch sms template item")

		return 0, err
	}

	// fetch subscription to get the pricing info
	var subscription models.Subscription
	err = repo.GetRepository().Find(ctx, &subscription, map[string]interface{}{
		constants.Type:    "SMS",
		constants.SubType: subCategory,
		constants.UserID:  userID,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return 0, err
	}

	var pricing models.Pricing
	err = repo.GetRepository().FindByID(ctx, subscription.PricingId, &pricing)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return 0, err
	}

	creditsToUse := getCreditsUseByLength(smsTemplate.Length)

	return creditsToUse * pricing.Rates, err
}

// getCreditsUseByLength gives us the credits consumption as per template length
// todo: update logic
func getCreditsUseByLength(templateLength int) float64 {

	if templateLength <= 160 {
		return CreditUseBelow160Length
	}

	return CreditUseAbove160Length
}
