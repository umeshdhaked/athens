package otp

import (
	"errors"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"gorm.io/gorm"
	"sync"
	"time"

	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/otp"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/gin-gonic/gin"
)

var (
	once    sync.Once
	service *OtpService
)

type OtpService struct {
	otpSender *otp.OtpSender
	crypto    *crypto.Crypto
	baseRepo  repo.IRepository
	otpRepo   repo.IOtpRepo
}

func NewOtpService(otpSender *otp.OtpSender, crypto *crypto.Crypto) {
	once.Do(func() {
		service = &OtpService{
			otpSender: otpSender,
			crypto:    crypto,
			baseRepo:  repo.GetRepository(),
			otpRepo:   repo.GetOtpRepo()}
	})
}

func GetOtpService() *OtpService {
	return service
}

func (o *OtpService) SendOtp(ctx *gin.Context, mobile string) error {
	generatedOtp := o.otpSender.GenerateOtp()
	logger.GetLogger().WithField("generatedOtp", generatedOtp).Info("generated otp:")
	if err := o.otpSender.SendOtp(generatedOtp); err != nil {
		return err
	}
	hashedOtp := o.crypto.HashString(generatedOtp)
	otp := models.Otp{
		Mobile: mobile,
		Otp:    hashedOtp,
		Exp:    time.Now().Add(2 * time.Minute).Unix(),
	}

	var fetchedOtp models.Otp
	err := o.baseRepo.Find(ctx, &fetchedOtp, map[string]interface{}{
		models.ColumnOtpMobile: mobile,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := o.baseRepo.Create(ctx, &otp); err != nil {
			return err
		}
		return nil
	} else {
		otp.ID = fetchedOtp.ID
		if err := o.baseRepo.Update(ctx, &otp); err != nil {
			return err
		}
		return nil
	}

}

func (o *OtpService) VerifyOtp(ctx *gin.Context, mobile string, otp string) error {
	currentHashedOtp := o.crypto.HashString(otp)

	var fetchedOtp models.Otp
	err := o.baseRepo.Find(ctx, &fetchedOtp, map[string]interface{}{
		models.ColumnOtpMobile: mobile,
	})
	if err != nil {
		return err
	}

	currTime := time.Now().Unix()
	if fetchedOtp.Otp == currentHashedOtp && fetchedOtp.Exp > currTime {
		return nil
	}

	return errors.New("otp verification failed")
}
