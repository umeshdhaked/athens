package subscription

import (
	"errors"
	"time"

	"github.com/umeshdhaked/athens/internal/constants"
	"github.com/umeshdhaked/athens/internal/models"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/umeshdhaked/athens/pkg/logger"
	"github.com/gin-gonic/gin"
	gormLogger "gorm.io/gorm/logger"
)

func (s *SubscriptionService) CreateNewPricingSystem(ctx *gin.Context, request *dtos.PricingRequest) (*dtos.PricingResponse, error) {

	if role, exists := ctx.Get(constants.JwtTokenRole); !exists {
		return nil, errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return nil, errors.New("only admin user allowed to add pricing")
	}

	if "" == request.Category || "" == request.SubCatgory || "" == request.Type || 0 == request.Rates {
		return nil, errors.New("invalid input")
	}
	if !(request.Category == "SMS" || request.Category == "EMAIL" || request.Category == "WHATSAPP") {
		return nil, errors.New("invalid category")
	}
	if request.Category == "SMS" && request.SubCatgory != "PROMOTIONAL" && request.SubCatgory != "TRANSACTIONAL" {
		return nil, errors.New("invalid sub_category for SMS")
	}

	// search in DB if default exists.
	if request.Type == "DEFAULT" {
		var pricing models.Pricing
		err := s.baseRepo.Find(ctx, &pricing, map[string]interface{}{
			constants.Category:    request.Category,
			constants.SubCategory: request.SubCatgory,
		})
		if err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound) {
			logger.GetLogger().Error(err.Error())
			return nil, err
		}

		if pricing.ID != 0 {
			logger.GetLogger().Error("default pricing already exists for category and subcategory")
			return nil, errors.New("default pricing already exists for category and subcategory")
		}
	}

	// save pricing to db.
	obj := models.Pricing{
		Category:    request.Category,
		SubCategory: request.SubCatgory,
		PricingType: request.Type,
		Rates:       request.Rates,
		State:       "ACTIVE",
		BaseModel:   models.BaseModel{CreatedAt: time.Now().Unix()},
	}

	if er := s.baseRepo.Create(ctx, &obj); er != nil {
		return nil, er
	}

	return &dtos.PricingResponse{Id: obj.ID, Category: obj.Category, SubCatgory: obj.SubCategory, Type: obj.PricingType, Rates: obj.Rates, Status: obj.State}, nil
}

func (s *SubscriptionService) FetchAllActivePricingModel(ctx *gin.Context) ([]*dtos.PricingResponse, error) {
	var pricing []models.Pricing
	er := s.baseRepo.FindMultiple(ctx, &pricing, map[string]interface{}{
		// todo: there is no status here.
	})
	if er != nil {
		return nil, er
	}

	resp := []*dtos.PricingResponse{}
	for _, p := range pricing {
		resp = append(resp, &dtos.PricingResponse{
			Id:         p.ID,
			Category:   p.Category,
			SubCatgory: p.SubCategory,
			Type:       p.PricingType,
			Rates:      p.Rates,
			Status:     p.State,
		})
	}

	return resp, nil
}

func (s *SubscriptionService) PricingStatusUpdate(ctx *gin.Context, pricingReq *dtos.DeactivatePricingRequest) error {
	if role, exists := ctx.Get(constants.JwtTokenRole); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	var pricing models.Pricing
	err := s.baseRepo.FindByID(ctx, pricingReq.Id, &pricing)
	if err != nil {
		return err
	}

	pricing.State = "INACTIVE"
	if pricingReq.Enable {
		pricing.State = "ACTIVE"
	}

	return s.baseRepo.Create(ctx, &pricing)
}
