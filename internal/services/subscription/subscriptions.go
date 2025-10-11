package subscription

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"

	"github.com/umeshdhaked/athens/internal/constants"
	"github.com/umeshdhaked/athens/internal/models"
	"github.com/umeshdhaked/athens/internal/pkg/repo"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/umeshdhaked/athens/pkg/logger"
	"github.com/gin-gonic/gin"
)

var (
	once    sync.Once
	service *SubscriptionService
)

type SubscriptionService struct {
	baseRepo        repo.IRepository
	pricingRepo     repo.IPricingRepo
	subRepo         repo.ISubscriptionRepo
	userRepo        repo.IUserRepo
	creditRepo      repo.ICreditsRepo
	creditAuditRepo repo.ICreditsAuditRepo
}

func NewSubscriptionService(pricingRepo repo.IPricingRepo, subRepo repo.ISubscriptionRepo, userRepo repo.IUserRepo, creditRepo repo.ICreditsRepo, creditAuditRepo repo.ICreditsAuditRepo) {
	once.Do(func() {
		service = &SubscriptionService{
			baseRepo:        repo.GetRepository(),
			pricingRepo:     pricingRepo,
			subRepo:         subRepo,
			userRepo:        userRepo,
			creditRepo:      creditRepo,
			creditAuditRepo: creditAuditRepo}
	})
}

func GetSubscriptionService() *SubscriptionService {
	return service
}

func (s *SubscriptionService) AddDefaultSubscriptionToUser(ctx *gin.Context, request *dtos.UserDefaultSubscriptionRequest) error {
	if role, exists := ctx.Get(constants.JwtTokenRole); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	// get User
	var user models.User
	err := s.baseRepo.Find(ctx, &user, map[string]interface{}{
		models.ColumnUserMobile: request.UserMobile,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	// get default Pricing models
	var defaultPricings []models.Pricing
	err = s.baseRepo.FindMultiple(ctx, &defaultPricings, map[string]interface{}{
		constants.PricingType: "DEFAULT",
		constants.State:       "ACTIVE",
	})
	if err != nil {
		return err
	}
	admin := ctx.GetInt64(constants.JwtTokenUserID)

	// create Subscriptions to USER
	userSubsDto := []models.Subscription{}
	for _, pricing := range defaultPricings {
		userSubsDto = append(userSubsDto, models.Subscription{
			PricingId: pricing.ID,
			UserId:    user.ID,
			Type:      pricing.Category,
			SubType:   pricing.SubCategory,
			Status:    "ACTIVE",
			AddedBy:   strconv.Itoa(int(admin)),
			BaseModel: models.BaseModel{CreatedAt: time.Now().Unix()},
		})
	}

	return s.baseRepo.CreateInBatches(ctx, userSubsDto, len(userSubsDto))
}

func (s *SubscriptionService) UpdateSubscriptionToUser(ctx *gin.Context, request *dtos.UserSubscriptionRequest) error {
	if role, exists := ctx.Get(constants.JwtTokenRole); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	// get User
	var user models.User
	err := s.baseRepo.Find(ctx, &user, map[string]interface{}{
		models.ColumnUserMobile: request.UserMobile,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	var pricing models.Pricing
	err = s.baseRepo.FindByID(ctx, request.PricingId, &pricing)
	if err != nil {
		return err
	}

	admin := ctx.GetInt64(constants.JwtTokenUserID)

	// update existing Subscriptions for type-subtype
	userSubsDto := models.Subscription{
		PricingId: pricing.ID,
		UserId:    user.ID,
		Type:      pricing.Category,
		SubType:   pricing.SubCategory,
		Status:    "ACTIVE",
		AddedBy:   strconv.Itoa(int(admin)),
		BaseModel: models.BaseModel{CreatedAt: time.Now().Unix()},
	}

	// get existing subscription from type and subtype
	var subscription models.Subscription
	err = s.baseRepo.Find(ctx, &subscription, map[string]interface{}{
		constants.UserID:  user.ID,
		constants.Type:    pricing.Category,
		constants.SubType: pricing.SubCategory,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.GetLogger().Error(err.Error())
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if er := s.baseRepo.Create(ctx, &userSubsDto); er != nil {
			return errors.Join(er, errors.New("unable to create subscription for user"))
		}
	} else {
		userSubsDto.ID = subscription.ID
		if er := s.baseRepo.Update(ctx, &userSubsDto); er != nil {
			return errors.Join(er, errors.New("unable to update subscription for user"))
		}
	}

	return nil
}

func (s *SubscriptionService) FetchAllActiveSubscriptionsForUser(ctx *gin.Context, request *dtos.FetchSubscriptionRequest) ([]*dtos.SubscriptionResponse, error) {
	// get User
	var user models.User
	err := s.baseRepo.Find(ctx, &user, map[string]interface{}{
		models.ColumnUserMobile: request.UserMobile,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	var subscriptionEntities []models.Subscription
	err = s.baseRepo.FindMultiple(ctx, &subscriptionEntities, map[string]interface{}{
		constants.UserID: user.ID,
		constants.Status: "ACTIVE",
	})
	if err != nil {
		return nil, err
	}

	var resp []*dtos.SubscriptionResponse
	for _, subscription := range subscriptionEntities {
		resp = append(resp, &dtos.SubscriptionResponse{
			Id:        subscription.ID,
			PricingId: subscription.PricingId,
			UserId:    subscription.UserId,
			Type:      subscription.Type,
			SubType:   subscription.SubType,
			Status:    subscription.Status,
			AddedBy:   subscription.AddedBy,
			CreatedAt: subscription.CreatedAt,
			DeletedAt: subscription.DeletedAt,
		})
	}
	return resp, nil
}

func (s *SubscriptionService) SubscriptionsStatusUpdate(ctx *gin.Context, subRequest *dtos.DeactivateSubscriptionRequest) error {
	var subsription models.Subscription
	err := s.baseRepo.FindByID(ctx, subRequest.Id, &subsription)
	if err != nil {
		return err
	}

	subsription.Status = "INACTIVE"
	if subRequest.Enable {
		subsription.Status = "ACTIVE"
	}

	// todo: why create?
	return s.baseRepo.Create(ctx, &subsription)
}
