package otp

import (
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"io"
	"log"

	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
	"github.com/fastbiztech/hastinapura/internal/pkg/repositories"
)

type OtpSender struct {
	otpRepo *repositories.OtpRepo
}

func NewOtpSender(otpRepo *repositories.OtpRepo) *OtpSender {
	return &OtpSender{otpRepo: otpRepo}
}

func (o *OtpSender) GenerateOtp() string {

	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func (o *OtpSender) SendOtp(otp string) error {
	log.Printf("otp sent: %s", otp)
	//send otp here
	return nil
}

func (o *OtpSender) SaveOtp(ctx *gin.Context, mobile string, hashedOtp string) error {
	if err := o.otpRepo.SaveOtp(ctx, mobile, hashedOtp); err != nil {
		return err
	}
	return nil
}

func (o *OtpSender) FetchOtp(ctx *gin.Context, mobileNo string) *dbo.Otp {
	if otp, err := o.otpRepo.GetOtp(ctx, mobileNo); err != nil {
		return nil
	} else {
		return otp
	}
}
