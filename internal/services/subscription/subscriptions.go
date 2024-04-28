package subscription

import (
	"errors"
	"sync"
	"time"

	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
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
			SubType:   dp.SubCategory,
			SubStatus: "INACTIVE",
			AddedBy:   admin.(string),
			BaseModel: models.BaseModel{CreatedAt: time.Now().Unix()},
		})
	}

	return s.subRepo.BatchCreateUserSubscription(ctx, userSubsDto)
}

func (s *SubscriptionService) UpdateSubscriptionToUser(ctx *gin.Context, subReq *dtos.UserSubscriptionRequest) error {
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

	admin, _ := ctx.Get(constants.JwtTokenUserID)

	// update existing Subscriptions for type-subtype
	userSubsDto := models.UserSubscription{
		Id:        uuid.New().String(),
		PricingId: pricing.Id,
		UserId:    user.ID,
		Type:      pricing.Category,
		SubType:   pricing.SubCategory,
		SubStatus: "ACTIVE",
		AddedBy:   admin.(string),
		BaseModel: models.BaseModel{CreatedAt: time.Now().Unix()},
	}

	// get existing subscription from type and subtype
	existingsSubs, err := s.subRepo.FetchSubscriptionByTypeSubType(ctx, user.ID, pricing.Category, pricing.SubCategory)
	if err != nil {
		return err
	}
	if len(existingsSubs) > 0 {
		userSubsDto = models.UserSubscription{
			Id: existingsSubs[0].Id,
		}
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

	subscriptions, err := s.subRepo.FetchAllSubscriptionByStatus(ctx, user.ID, "ACTIVE")
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
			Status:    s.SubStatus,
			AddedBy:   s.AddedBy,
			CreatedAt: s.CreatedAt,
			DeletedAt: s.DeletedAt,
		})
	}
	return resp, nil
}

func (s *SubscriptionService) SubscriptionsStatusUpdate(ctx *gin.Context, subRequest *dtos.DeactivateSubscriptionRequest) error {
	subsription, err := s.subRepo.GetSubscriptionFromId(ctx, subRequest.Id)
	if err != nil {
		return err
	}

	subsription.SubStatus = "INACTIVE"
	if subRequest.Enable {
		subsription.SubStatus = "ACTIVE"
	}

	return s.subRepo.CreateUserSubscription(ctx, subsription)
}
