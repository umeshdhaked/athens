package promo

import (
	"errors"
	"github.com/umeshdhaked/athens/pkg/logger"
	"gorm.io/gorm"
	"sync"

	"github.com/umeshdhaked/athens/internal/models"
	"github.com/gin-gonic/gin"

	"github.com/umeshdhaked/athens/internal/pkg/repo"
)

var (
	once    sync.Once
	service *PromoService
)

type PromoService struct {
	baseRepo  repo.IRepository
	promoRepo repo.IPromotionRepo
}

func NewPromoService(promoRepo repo.IPromotionRepo) {
	once.Do(func() {
		service = &PromoService{
			baseRepo:  repo.GetRepository(),
			promoRepo: promoRepo,
		}
	})
}

func GetPromoService() *PromoService {
	return service
}

func (s *PromoService) SavePhoneNo(ctx *gin.Context, phoneNo string) error {
	var exPromoPh models.PromoPhone
	err := s.baseRepo.Find(ctx, &exPromoPh, map[string]interface{}{
		models.ColumnPromoPhoneMobile: phoneNo,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	obj := models.PromoPhone{
		Mobile: phoneNo,
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) && exPromoPh.IsAlreadyContacted == "true" {
		logger.GetLogger().Info("promo phone existed and already contacted")
		return nil
	} else {
		obj.IsAlreadyContacted = "false"
	}

	return s.baseRepo.Create(ctx, &obj)
}

func (s *PromoService) FetchPromoNumbers(ctx *gin.Context, isAlreadyConnected string) ([]models.PromoPhone, error) {
	var items []models.PromoPhone
	err := s.baseRepo.FindMultiple(ctx, &items, map[string]interface{}{
		models.ColumnIsAlreadyContacted: isAlreadyConnected,
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *PromoService) MarkContacted(ctx *gin.Context, mobile string, comment string) error {
	var exPromoPh models.PromoPhone
	err := s.baseRepo.Find(ctx, &exPromoPh, map[string]interface{}{
		models.ColumnPromoPhoneMobile: mobile,
	})
	if err != nil {
		return err
	}

	obj := models.PromoPhone{
		ID:                 exPromoPh.ID,
		Mobile:             mobile,
		IsAlreadyContacted: "true",
		Comment:            comment,
	}

	return s.baseRepo.Update(ctx, &obj)
}
