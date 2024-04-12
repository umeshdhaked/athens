package otp

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/otp"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	once    sync.Once
	service *OtpService
)

type OtpService struct {
	otpSender *otp.OtpSender
	crypto    *crypto.Crypto
	otpRepo   *repo.OtpRepo
}

func NewOtpService(otpSender *otp.OtpSender, crypto *crypto.Crypto, otpRepo *repo.OtpRepo) {
	once.Do(func() {
		service = &OtpService{otpSender: otpSender, crypto: crypto, otpRepo: otpRepo}
	})
}

func GetOtpService() *OtpService {
	return service
}

func (o *OtpService) SendOtp(ctx *gin.Context, mobile string) error {
	generatedOtp := o.otpSender.GenerateOtp()
	log.Printf("generated otp %s", generatedOtp)
	if err := o.otpSender.SendOtp(generatedOtp); err != nil {
		return err
	}
	hashedOtp := o.crypto.HashString(generatedOtp)

	otp := models.Otp{
		Id:     uuid.New().String(),
		Mobile: mobile,
		Otp:    hashedOtp,
		Exp:    time.Now().Add(2 * time.Minute).Unix(),
	}
	if err := o.otpRepo.SaveOtp(ctx, otp); err != nil {
		return err
	}

	return nil
}

func (o *OtpService) VerifyOtp(ctx *gin.Context, mobile string, otp string) error {
	currentHashedOtp := o.crypto.HashString(otp)
	fetchedOtp, err := o.otpRepo.GetOtp(ctx, mobile)
	if err != nil {
		return err
	}

	currTime := time.Now().Unix()
	if fetchedOtp != nil && fetchedOtp.Otp == currentHashedOtp && fetchedOtp.Exp > currTime {
		return nil
	}
	return errors.New("otp verification failed")
}
