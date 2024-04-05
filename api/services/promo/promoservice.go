package promo

import (
	"log"
	"time"

	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
	"github.com/fastbiztech/hastinapura/internal/pkg/repositories"
)

type PromoService struct {
	promoRepo *repositories.PromotionRepo
}

func NewPromoService(promoRepo *repositories.PromotionRepo) *PromoService {
	return &PromoService{promoRepo: promoRepo}
}

func (s *PromoService) SavePhoneNo(phoneNo string) error {
	exPromoPh, err := s.promoRepo.GetPromoFromMobile(phoneNo)
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

	return s.promoRepo.AddPromoContact(&obj)
}

func (s *PromoService) FetchPromoNumbers(isAlreadyConnected string) ([]dbo.PromoPhone, error) {
	return s.promoRepo.GetAlreadyContactedPromo(isAlreadyConnected)
}

func (s *PromoService) MarkContacted(mobile string, comment string) error {
	obj := dbo.PromoPhone{Mobile: mobile,
		Timestamp: time.Now().Format(time.RFC850), IsAlreadyContacted: "true", Comment: comment}

	return s.promoRepo.MarkContacted(&obj)
}
