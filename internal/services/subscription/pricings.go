package subscription

import (
	"errors"
	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

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
		SubCategory:  pricing.SubCatgory,
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
			SubCatgory: p.SubCategory,
			Type:       p.PricingType,
			Rates:      p.Rates,
			Status:     p.PricingState,
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

	pricing, err := s.pricingRepo.GetPricingByPricingID(ctx, pricingReq.Id)
	if err != nil {
		return err
	}

	pricing.PricingState = "INACTIVE"
	if pricingReq.Enable {
		pricing.PricingState = "ACTIVE"
	}

	return s.pricingRepo.CreatePricing(ctx, pricing)
}
