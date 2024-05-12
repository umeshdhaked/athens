package subscription

import (
	"errors"
	"log"
	"sort"
	"time"

	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *SubscriptionService) AddCreditToUser(ctx *gin.Context, subRequest *dtos.AddCreditsRequest) error {
	user, err := s.userRepo.GetUserFromMobile(ctx, subRequest.UserMobile)
	if err != nil {
		return err
	}

	subscriptions, err := s.subRepo.FetchAllSubscriptionByStatus(ctx, user.ID, "INACTIVE")
	if err != nil {
		return err
	}

	var userSubsDto []models.Subscription
	if nil != subscriptions {
		for _, sub := range subscriptions {
			sub.Status = "ACTIVE"
			userSubsDto = append(userSubsDto, sub)
		}
		err = s.subRepo.BatchCreateUserSubscription(ctx, userSubsDto)
		if err != nil {
			return err
		}
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
		Id:             uuid.New().String(),
		DeductedAmount: 0,
		AddedAmount:    subRequest.InitialCredit,
		CreditId:       credit.ID,
		UserId:         credit.UserID,
		PaymentOrderId: subRequest.PaymentOrderId,
		BaseModel:      models.BaseModel{UpdatedAt: time.Now().Unix()},
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
	usersSubscription, er := s.subRepo.FetchAllSubscriptionByStatus(ctx, userId, "ACTIVE")
	if er != nil {
		return er
	}
	sort.Slice(usersSubscription[:], func(i, j int) bool {
		return usersSubscription[i].CreatedAt > usersSubscription[j].CreatedAt
	})
	var currentActiveSub models.Subscription
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
		Id:             uuid.New().String(),
		Category:       category,
		SubCategory:    subCategory,
		DeductedAmount: totalUsedCredit,
		AddedAmount:    0,
		CreditId:       credit.ID,
		UserId:         credit.UserID,
		BaseModel:      models.BaseModel{UpdatedAt: time.Now().Unix()},
	}); err != nil {
		return err
	}

	return s.creditRepo.CreateUserCredit(ctx, credit)
}
