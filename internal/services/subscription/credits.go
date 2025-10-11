package subscription

import (
	"errors"
	"sort"
	"time"

	"github.com/umeshdhaked/athens/internal/constants"
	"github.com/umeshdhaked/athens/internal/models"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/gin-gonic/gin"
	gormLogger "gorm.io/gorm/logger"
)

func (s *SubscriptionService) AddCreditToUser(ctx *gin.Context, subRequest *dtos.AddCreditsRequest) error {
	var user models.User
	err := s.baseRepo.Find(ctx, &user, map[string]interface{}{
		models.ColumnUserMobile: subRequest.UserMobile,
	})
	if err != nil {
		return err
	}

	var subscriptions []models.Subscription
	err = s.baseRepo.FindMultiple(ctx, &subscriptions, map[string]interface{}{
		constants.UserID: user.ID,
		constants.Status: "INACTIVE",
	})
	if err != nil {
		return err
	}

	var userSubsDto []models.Subscription
	if nil != subscriptions && len(subscriptions) > 0 {
		for _, sub := range subscriptions {
			sub.Status = "ACTIVE"
			userSubsDto = append(userSubsDto, sub)
		}
		err = s.baseRepo.CreateInBatches(ctx, userSubsDto, len(userSubsDto))
		if err != nil {
			return err
		}
	}

	credit := &models.Credits{}
	er := s.baseRepo.Find(ctx, credit, map[string]interface{}{
		constants.UserID: user.ID,
	})
	if er != nil && !errors.Is(er, gormLogger.ErrRecordNotFound) {
		return err
	}

	// Do credit audit
	if err = s.baseRepo.Create(ctx, &models.CreditAudits{
		DeductedAmount: 0,
		AddedAmount:    subRequest.InitialCredit,
		CreditId:       credit.ID,
		UserId:         credit.UserId,
		PaymentOrderId: subRequest.PaymentOrderId,
		BaseModel:      models.BaseModel{UpdatedAt: time.Now().Unix()},
	}); err != nil {
		return err
	}

	if errors.Is(er, gormLogger.ErrRecordNotFound) { // if not found then create
		credit = &models.Credits{
			UserId:      user.ID,
			Balance:     subRequest.InitialCredit,
			BalanceLeft: subRequest.InitialCredit,
		}
		return s.baseRepo.Create(ctx, credit)
	} else {
		credit.Balance = credit.Balance + subRequest.InitialCredit
		credit.BalanceLeft = credit.BalanceLeft + subRequest.InitialCredit
		return s.baseRepo.Update(ctx, credit)
	}

}

func (s *SubscriptionService) FetchCredit(ctx *gin.Context) (*dtos.CreditsResponse, error) {
	userId := ctx.GetInt64(constants.JwtTokenUserID)
	if 0 == userId { // create another version of this with payment validation with transaction ID
		return nil, errors.New("internal server error, user id not found")
	}
	mobile, exists := ctx.Get(constants.JwtTokenMobile)
	if !exists { // create another version of this with payment validation with transaction ID
		return nil, errors.New("internal server error, user mobile not found")
	}

	var credit models.Credits
	err := s.baseRepo.Find(ctx, &credit, map[string]interface{}{
		constants.UserID: userId,
	})
	if err != nil {
		return nil, err
	}

	creditResp := &dtos.CreditsResponse{
		Id:              credit.ID,
		UserMobile:      mobile.(string),
		InitialCredit:   credit.Balance,
		RemainingCredit: credit.BalanceLeft,
		CreatedAt:       credit.CreatedAt,
	}

	return creditResp, nil
}

func (s *SubscriptionService) ChargeUser(ctx *gin.Context, userId string, category string, subCategory string, unitCount float64) error {
	// get subscription (recent one for typ,subtype,userId)
	var usersSubscription []models.Subscription
	er := s.baseRepo.FindMultiple(ctx, &usersSubscription, map[string]interface{}{
		constants.Status: "ACTIVE",
	})
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
	var pricing models.Pricing
	err := s.baseRepo.FindByID(ctx, pricingId, &pricing)
	if err != nil {
		return er
	}

	// fetch credit for user_id
	var credit models.Credits
	err = s.baseRepo.Find(ctx, &credit, map[string]interface{}{
		constants.UserID: userId,
	})
	if err != nil {
		return err
	}

	totalUsedCredit := pricing.Rates * unitCount
	remainingCredit := credit.BalanceLeft

	// check and update credit of user
	if remainingCredit < totalUsedCredit {
		return errors.New("CREDIT_EXHAUSTED")
	}

	credit.BalanceLeft = credit.BalanceLeft - totalUsedCredit

	if err = s.baseRepo.Create(ctx, &models.CreditAudits{
		Category:       category,
		SubCategory:    subCategory,
		DeductedAmount: totalUsedCredit,
		AddedAmount:    0,
		CreditId:       credit.ID,
		UserId:         credit.UserId,
		BaseModel:      models.BaseModel{UpdatedAt: time.Now().Unix()},
	}); err != nil {
		return err
	}

	//todo: why are we creating it again?
	return s.baseRepo.Create(ctx, &credit)
}
