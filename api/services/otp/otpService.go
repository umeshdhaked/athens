package otp

import (
	"errors"
	"log"

	"github.com/FastBizTech/hastinapura/pkg/services/crypto"
	"github.com/FastBizTech/hastinapura/pkg/services/otp"
)

type OtpService struct {
	otpSender *otp.OtpSender
	crypto    *crypto.Crypto
}

func NewOtpService(otpSender *otp.OtpSender, crypto *crypto.Crypto) *OtpService {
	return &OtpService{otpSender: otpSender, crypto: crypto}
}

func (o *OtpService) SendOtp(mobile string) error {
	generatedOtp := o.otpSender.GenerateOtp()
	log.Printf("generated otp %s", generatedOtp)
	if err := o.otpSender.SendOtp(generatedOtp); err != nil {
		return err
	}
	hashedOtp := o.crypto.HashString(generatedOtp)
	if err := o.otpSender.SaveOtp(mobile, hashedOtp); err != nil {
		return err
	}
	return nil
}

func (o *OtpService) VerifyOtp(mobile string, otp string) error {
	currentHashedOtp := o.crypto.HashString(otp)
	hashedOtp := o.otpSender.FetchOtp(mobile)

	if hashedOtp == currentHashedOtp {
		return nil
	}
	return errors.New("otp verification failed")
}
