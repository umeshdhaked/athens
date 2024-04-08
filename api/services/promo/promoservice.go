package promo

import (
	"log"
	"time"

	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/gin-gonic/gin"

	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
)

type PromoService struct {
	promoRepo *repo.PromotionRepo
}

func NewPromoService(promoRepo *repo.PromotionRepo) *PromoService {
	return &PromoService{promoRepo: promoRepo}
}

func (s *PromoService) SavePhoneNo(ctx *gin.Context, phoneNo string) error {
	exPromoPh, err := s.promoRepo.GetPromoFromMobile(ctx, phoneNo)
	if err != nil {
		return err
	}

	obj := dbo.PromoPhone{Mobile: phoneNo, Timestamp: time.Now().Format(time.RFC850)}
	if exPromoPh != nil && exPromoPh.IsAlreadyContacted == "true" {
		log.Print("promo phone existed and already contacted")
		return nil
	} else {
		obj.IsAlreadyContacted = "false"
	}

	return s.promoRepo.AddPromoContact(ctx, &obj)
}

func (s *PromoService) FetchPromoNumbers(ctx *gin.Context, isAlreadyConnected string) ([]dbo.PromoPhone, error) {
	return s.promoRepo.GetAlreadyContactedPromo(ctx, isAlreadyConnected)
}

func (s *PromoService) MarkContacted(ctx *gin.Context, mobile string, comment string) error {
	obj := dbo.PromoPhone{Mobile: mobile,
		Timestamp: time.Now().Format(time.RFC850), IsAlreadyContacted: "true", Comment: comment}

	return s.promoRepo.MarkContacted(ctx, &obj)
}
