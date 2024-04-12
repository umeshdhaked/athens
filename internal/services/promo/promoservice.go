package promo

import (
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"time"

	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
)

var (
	once    sync.Once
	service *PromoService
)

type PromoService struct {
	promoRepo *repo.PromotionRepo
}

func NewPromoService(promoRepo *repo.PromotionRepo) {
	once.Do(func() {
		service = &PromoService{promoRepo: promoRepo}
	})
}

func GetPromoService() *PromoService {
	return service
}

func (s *PromoService) SavePhoneNo(ctx *gin.Context, phoneNo string) error {
	exPromoPh, err := s.promoRepo.GetPromoFromMobile(ctx, phoneNo)
	if err != nil {
		return err
	}

	obj := models.PromoPhone{Mobile: phoneNo, Timestamp: time.Now().Format(time.RFC850)}
	if exPromoPh != nil && exPromoPh.IsAlreadyContacted == "true" {
		log.Print("promo phone existed and already contacted")
		return nil
	} else {
		obj.IsAlreadyContacted = "false"
	}

	return s.promoRepo.AddPromoContact(ctx, &obj)
}

func (s *PromoService) FetchPromoNumbers(ctx *gin.Context, isAlreadyConnected string) ([]models.PromoPhone, error) {
	return s.promoRepo.GetAlreadyContactedPromo(ctx, isAlreadyConnected)
}

func (s *PromoService) MarkContacted(ctx *gin.Context, mobile string, comment string) error {
	obj := models.PromoPhone{Mobile: mobile,
		Timestamp: time.Now().Format(time.RFC850), IsAlreadyContacted: "true", Comment: comment}

	return s.promoRepo.MarkContacted(ctx, &obj)
}
