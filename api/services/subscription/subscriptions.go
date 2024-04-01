package subscription

import (
	"errors"
	"time"

	"github.com/fastbiztech/hastinapura/pkg/models/dbo"
	"github.com/fastbiztech/hastinapura/pkg/models/requests"
	"github.com/fastbiztech/hastinapura/pkg/models/responses"
	"github.com/fastbiztech/hastinapura/pkg/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	pricingRepo *repositories.PricingRepo
	subRepo     *repositories.SubscriptionRepo
	userRepo    *repositories.UserRepo
}

func NewSubscriptionService(pricingRepo *repositories.PricingRepo, subRepo *repositories.SubscriptionRepo, userRepo *repositories.UserRepo) *SubscriptionService {
	return &SubscriptionService{pricingRepo: pricingRepo, subRepo: subRepo, userRepo: userRepo}
}

func (s *SubscriptionService) CreateNewPricingSystem(ctx *gin.Context, pricing *requests.PricingRequest) (*responses.PricingResponse, error) {

	if role, exists := ctx.Params.Get("role"); !exists {
		return nil, errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return nil, errors.New("only admin user allowed to add pricing")
	}

	if "" == pricing.Category || "" == pricing.SubCatgory || "" == pricing.Type || 0 == pricing.Rates {
		return nil, errors.New("invalid input")
	}
	if !(pricing.Category == "SMS" || pricing.Category == "EMAIL" || pricing.Category == "WHATSAPP") {
		return nil, errors.New("inavaid category")
	}
	if pricing.Category == "SMS" && pricing.SubCatgory != "PROMOTIONAL" && pricing.SubCatgory != "TRANSACTIONAL" {
		return nil, errors.New("inavaid sub_category for SMS")
	}

	// search in DB if default exists.
	if pricing.Type == "DEFAULT" {
		resp, err := s.pricingRepo.GetDefaultPricingsForCategoryAndSubCategory(pricing.Category, pricing.SubCatgory)
		if err != nil {
			return nil, err
		}
		if resp != nil && len(resp) > 0 {
			return nil, errors.New("default pricing already exists for category and subcategory")
		}
	}

	// save pricing to db.
	obj := dbo.Pricing{
		Id:           uuid.New().String(),
		Category:     pricing.Category,
		SubCatgory:   pricing.SubCatgory,
		PricingType:  pricing.Type,
		Rates:        pricing.Rates,
		PricingState: "ACTIVE",
		CreatedAt:    time.Now().Unix(),
	}

	if er := s.pricingRepo.CreatePricing(&obj); er != nil {
		return nil, er
	}

	return &responses.PricingResponse{Id: obj.Id}, nil
}

func (s *SubscriptionService) FetchAllActivePricingModel(ctx *gin.Context) ([]*responses.PricingResponse, error) {
	if role, exists := ctx.Params.Get("role"); !exists {
		return nil, errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return nil, errors.New("only admin user allowed to add pricing")
	}

	pricing, er := s.pricingRepo.FetchAllActivePricing()
	if er != nil {
		return nil, er
	}

	resp := []*responses.PricingResponse{}
	for _, p := range pricing {
		resp = append(resp, &responses.PricingResponse{
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

func (s *SubscriptionService) AddDefaultSubscriptionToUser(ctx *gin.Context, subReq *requests.UserDefaultSubscriptionRequest) error {
	if role, exists := ctx.Params.Get("role"); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	// get User
	mobile := subReq.UserMobile
	user, er := s.userRepo.GetUserFromMobile(mobile)
	if er != nil {
		return er
	}

	// get default Pricing models
	defaultPricings, err := s.pricingRepo.GetAllDefaultActivePricings()
	if err != nil {
		return err
	}
	admin, _ := ctx.Params.Get("id")

	// create Subscriptions to USER
	userSubsDto := []dbo.UserSubscription{}
	for _, dp := range defaultPricings {
		userSubsDto = append(userSubsDto, dbo.UserSubscription{
			Id:        uuid.New().String(),
			PricingId: dp.Id,
			UserId:    user.Id,
			Type:      dp.Category,
			SubType:   dp.SubCatgory,
			Status:    "ACTIVE",
			AddedBy:   admin,
			CreatedAt: time.Now().Unix(),
		})
	}

	return s.subRepo.BatchCreateUserSubscription(ctx, userSubsDto)
}

func (s *SubscriptionService) AddSubscriptionToUser(ctx *gin.Context, subReq *requests.UserSubscriptionRequest) error {
	if role, exists := ctx.Params.Get("role"); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	// get User
	mobile := subReq.UserMobile
	user, er := s.userRepo.GetUserFromMobile(mobile)
	if er != nil {
		return er
	}

	pricing, err := s.pricingRepo.GetPricingByPricingID(subReq.PricingId)
	if err != nil {
		return err
	}

	// create Subscriptions to USER
	admin, _ := ctx.Params.Get("id")

	userSubsDto := dbo.UserSubscription{
		Id:        uuid.New().String(),
		PricingId: pricing.Id,
		UserId:    user.Id,
		Type:      pricing.Category,
		SubType:   pricing.SubCatgory,
		Status:    "ACTIVE",
		AddedBy:   admin,
		CreatedAt: time.Now().Unix(),
	}

	if er := s.subRepo.CreateUserSubscription(&userSubsDto); er != nil {
		return errors.Join(er, errors.New("unable to create subscription for user"))
	}
	return nil
}

func (s *SubscriptionService) FetchAllActiveSubscriptionsForUser(ctx *gin.Context, req *requests.FetchSubscriptionRequest) ([]*responses.SubscriptionResponse, error) {
	if role, exists := ctx.Params.Get("role"); !exists {
		return nil, errors.New("internal server error, user role not found")
	} else if role != "admin" {
		return nil, errors.New("only admin user allowed to add pricing")
	}

	// get User
	user, er := s.userRepo.GetUserFromMobile(req.UserMobile)
	if er != nil {
		return nil, er
	}

	subscriptions, err := s.subRepo.FetchAllSubscriptionForAUser(user.Id)
	if err != nil {
		return nil, err
	}

	resp := []*responses.SubscriptionResponse{}
	for _, s := range subscriptions {
		resp = append(resp, &responses.SubscriptionResponse{
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

// TODO: add/update credit api

func chargeUser(userId string, typ string, subtype string) {
	// get subscription (recent one for typ,subtype,userId)
	// get pricing from subscription
	// fetch credit for user_id
	// update credit from user_id
}
