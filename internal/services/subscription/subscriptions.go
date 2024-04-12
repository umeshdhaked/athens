package subscription

import (
	"errors"
	"github.com/fastbiztech/hastinapura/internal/models"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	once    sync.Once
	service *SubscriptionService
)

type SubscriptionService struct {
	pricingRepo     *repo.PricingRepo
	subRepo         *repo.SubscriptionRepo
	userRepo        *repo.UserRepo
	creditRepo      *repo.CreditsRepo
	creditAuditRepo *repo.CreditsAuditRepo
}

func NewSubscriptionService(pricingRepo *repo.PricingRepo, subRepo *repo.SubscriptionRepo, userRepo *repo.UserRepo, creditRepo *repo.CreditsRepo, creditAuditRepo *repo.CreditsAuditRepo) {
	once.Do(func() {
		service = &SubscriptionService{pricingRepo: pricingRepo, subRepo: subRepo, userRepo: userRepo, creditRepo: creditRepo, creditAuditRepo: creditAuditRepo}
	})
}

func GetSubscriptionService() *SubscriptionService {
	return service
}

func (s *SubscriptionService) CreateNewPricingSystem(ctx *gin.Context, pricing *dtos.PricingRequest) (*dtos.PricingResponse, error) {

	if role, exists := ctx.Get(constants.JwtTokenRole); !exists {
		return nil, errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return nil, errors.New("only admin user allowed to add pricing")
	}

	if "" == pricing.Category || "" == pricing.SubCatgory || "" == pricing.Type || 0 == pricing.Rates {
		return nil, errors.New("invalid input")
	}
	if !(pricing.Category == "SMS" || pricing.Category == "EMAIL" || pricing.Category == "WHATSAPP") {
		return nil, errors.New("invalid category")
	}
	if pricing.Category == "SMS" && pricing.SubCatgory != "PROMOTIONAL" && pricing.SubCatgory != "TRANSACTIONAL" {
		return nil, errors.New("invalid sub_category for SMS")
	}

	// search in DB if default exists.
	if pricing.Type == "DEFAULT" {
		resp, err := s.pricingRepo.GetDefaultPricingsForCategoryAndSubCategory(ctx, pricing.Category, pricing.SubCatgory)
		if err != nil {
			return nil, err
		}
		if resp != nil && len(resp) > 0 {
			return nil, errors.New("default pricing already exists for category and subcategory")
		}
	}

	// save pricing to db.
	obj := models.Pricing{
		Id:           uuid.New().String(),
		Category:     pricing.Category,
		SubCatgory:   pricing.SubCatgory,
		PricingType:  pricing.Type,
		Rates:        pricing.Rates,
		PricingState: "ACTIVE",
		BaseModel:    models.BaseModel{CreatedAt: time.Now().Unix()},
	}

	if er := s.pricingRepo.CreatePricing(ctx, &obj); er != nil {
		return nil, er
	}

	return &dtos.PricingResponse{Id: obj.Id}, nil
}

func (s *SubscriptionService) FetchAllActivePricingModel(ctx *gin.Context) ([]*dtos.PricingResponse, error) {
	pricing, er := s.pricingRepo.FetchAllActivePricing(ctx)
	if er != nil {
		return nil, er
	}

	resp := []*dtos.PricingResponse{}
	for _, p := range pricing {
		resp = append(resp, &dtos.PricingResponse{
			Id:         p.Id,
			Category:   p.Category,
			SubCatgory: p.SubCatgory,
			Type:       p.PricingType,
			Rates:      p.Rates,
			Status:     p.PricingState,
		})
	}

	return resp, nil
}

func (s *SubscriptionService) DeactivatePricing(ctx *gin.Context, pricingReq *dtos.DeactivatePricingRequest) error {
	if role, exists := ctx.Get(constants.JwtTokenRole); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	pricing, err := s.pricingRepo.GetPricingByPricingID(ctx, pricingReq.Id)
	if err != nil {
		return err
	}

	pricing.PricingState = "INACTIVE"

	return s.pricingRepo.CreatePricing(ctx, pricing)
}

func (s *SubscriptionService) AddDefaultSubscriptionToUser(ctx *gin.Context, subReq *dtos.UserDefaultSubscriptionRequest) error {
	if role, exists := ctx.Get(constants.JwtTokenRole); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	// get User
	mobile := subReq.UserMobile
	user, er := s.userRepo.GetUserFromMobile(ctx, mobile)
	if er != nil {
		return er
	}

	// get default Pricing models
	defaultPricings, err := s.pricingRepo.GetAllDefaultActivePricings(ctx)
	if err != nil {
		return err
	}
	admin, _ := ctx.Get(constants.JwtTokenUserID)

	// create Subscriptions to USER
	userSubsDto := []models.UserSubscription{}
	for _, dp := range defaultPricings {
		userSubsDto = append(userSubsDto, models.UserSubscription{
			Id:        uuid.New().String(),
			PricingId: dp.Id,
			UserId:    user.ID,
			Type:      dp.Category,
			SubType:   dp.SubCatgory,
			Status:    "ACTIVE",
			AddedBy:   admin.(string),
			BaseModel: models.BaseModel{CreatedAt: time.Now().Unix()},
		})
	}

	return s.subRepo.BatchCreateUserSubscription(ctx, userSubsDto)
}

func (s *SubscriptionService) AddSubscriptionToUser(ctx *gin.Context, subReq *dtos.UserSubscriptionRequest) error {
	if role, exists := ctx.Get(constants.JwtTokenRole); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	// get User
	mobile := subReq.UserMobile
	user, er := s.userRepo.GetUserFromMobile(ctx, mobile)
	if er != nil {
		return er
	}

	pricing, err := s.pricingRepo.GetPricingByPricingID(ctx, subReq.PricingId)
	if err != nil {
		return err
	}

	// create Subscriptions to USER
	admin, _ := ctx.Get(constants.JwtTokenUserID)

	userSubsDto := models.UserSubscription{
		Id:        uuid.New().String(),
		PricingId: pricing.Id,
		UserId:    user.ID,
		Type:      pricing.Category,
		SubType:   pricing.SubCatgory,
		Status:    "ACTIVE",
		AddedBy:   admin.(string),
		BaseModel: models.BaseModel{CreatedAt: time.Now().Unix()},
	}

	if er := s.subRepo.CreateUserSubscription(ctx, &userSubsDto); er != nil {
		return errors.Join(er, errors.New("unable to create subscription for user"))
	}
	return nil
}

func (s *SubscriptionService) FetchAllActiveSubscriptionsForUser(ctx *gin.Context, req *dtos.FetchSubscriptionRequest) ([]*dtos.SubscriptionResponse, error) {
	// get User
	user, er := s.userRepo.GetUserFromMobile(ctx, req.UserMobile)
	if er != nil {
		return nil, er
	}

	subscriptions, err := s.subRepo.FetchAllSubscriptionForAUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	resp := []*dtos.SubscriptionResponse{}
	for _, s := range subscriptions {
		resp = append(resp, &dtos.SubscriptionResponse{
			Id:        s.Id,
			PricingId: s.PricingId,
			UserId:    s.UserId,
			Type:      s.Type,
			SubType:   s.SubType,
			Status:    s.Status,
			AddedBy:   s.AddedBy,
			CreatedAt: s.CreatedAt,
			DeletedAt: s.DeletedAt,
		})
	}
	return resp, nil
}

func (s *SubscriptionService) DeactivateSubscriptionsForUser(ctx *gin.Context, subRequest *dtos.DeactivateSubscriptionRequest) error {
	subsription, err := s.subRepo.GetSubscriptionFromId(ctx, subRequest.Id)
	if err != nil {
		return err
	}

	subsription.Status = "INACTIVE"

	return s.subRepo.CreateUserSubscription(ctx, subsription)
}

func (s *SubscriptionService) AddCreditToUser(ctx *gin.Context, subRequest *dtos.AddCreditsRequest) error {
	user, err := s.userRepo.GetUserFromMobile(ctx, subRequest.UserMobile)
	if err != nil {
		return err
	}

	credit, err := s.creditRepo.FetchCreditByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	log.Println("credit added", subRequest)

	if credit == nil {
		credit = &models.Credits{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Credits:     subRequest.InitialCredit,
			CreditsLeft: subRequest.InitialCredit,
		}
	} else {
		credit.Credits = credit.Credits + subRequest.InitialCredit
		credit.CreditsLeft = credit.CreditsLeft + subRequest.InitialCredit
	}

	if err = s.creditAuditRepo.CreateUserCreditAudit(ctx, &models.CreditAudits{
		Id:            uuid.New().String(),
		DeductedAmout: 0,
		AddedAmount:   subRequest.InitialCredit,
		CreditId:      credit.ID,
		UserId:        credit.UserID,
		BaseModel:     models.BaseModel{UpdatedAt: time.Now().Unix()},
	}); err != nil {
		return err
	}

	return s.creditRepo.CreateUserCredit(ctx, credit)
}

func (s *SubscriptionService) FetchCredit(ctx *gin.Context) (*dtos.CreditsResponse, error) {
	userId, exists := ctx.Get(constants.JwtTokenUserID)
	if !exists { // create another version of this with payment validation with transaction ID
		return nil, errors.New("internal server error, user id not found")
	}
	mobile, exists := ctx.Get(constants.JwtTokenMobile)
	if !exists { // create another version of this with payment validation with transaction ID
		return nil, errors.New("internal server error, user mobile not found")
	}

	credit, err := s.creditRepo.FetchCreditByUserID(ctx, userId.(string))
	if err != nil {
		return nil, err
	}

	creditResp := &dtos.CreditsResponse{
		Id:              credit.ID,
		UserMobile:      mobile.(string),
		InitialCredit:   credit.Credits,
		RemainingCredit: credit.CreditsLeft,
		CreatedAt:       credit.CreatedAt,
	}

	return creditResp, nil
}

func (s *SubscriptionService) ChargeUser(ctx *gin.Context, userId string, category string, subCategory string, unitCount float64) error {
	// get subscription (recent one for typ,subtype,userId)
	usersSubscription, er := s.subRepo.FetchAllSubscriptionForAUser(ctx, userId)
	if er != nil {
		return er
	}
	sort.Slice(usersSubscription[:], func(i, j int) bool {
		return usersSubscription[i].CreatedAt > usersSubscription[j].CreatedAt
	})
	var currentActiveSub models.UserSubscription
	for _, sub := range usersSubscription {
		if sub.Status == "ACTIVE" && sub.Type == category && sub.SubType == subCategory {
			currentActiveSub = sub
			break
		}
	}

	pricingId := currentActiveSub.PricingId

	// get pricing from subscription
	pricing, err := s.pricingRepo.GetPricingByPricingID(ctx, pricingId)
	if err != nil {
		return er
	}

	// fetch credit for user_id
	credit, err := s.creditRepo.FetchCreditByUserID(ctx, userId)
	if err != nil {
		return err
	}
	if credit == nil {
		return errors.New("credit not found")
	}

	totalUsedCredit := pricing.Rates * unitCount
	remainingCredit := credit.CreditsLeft

	// check and update credit of user
	if remainingCredit < totalUsedCredit {
		return errors.New("CREDIT_EXHAUSTED")
	}

	credit.CreditsLeft = credit.CreditsLeft - totalUsedCredit

	if err = s.creditAuditRepo.CreateUserCreditAudit(ctx, &models.CreditAudits{
		Id:            uuid.New().String(),
		Category:      category,
		SubCategory:   subCategory,
		DeductedAmout: totalUsedCredit,
		AddedAmount:   0,
		CreditId:      credit.ID,
		UserId:        credit.UserID,
		BaseModel:     models.BaseModel{UpdatedAt: time.Now().Unix()},
	}); err != nil {
		return err
	}

	return s.creditRepo.CreateUserCredit(ctx, credit)
}
