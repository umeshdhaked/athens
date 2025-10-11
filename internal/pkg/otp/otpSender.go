package otp

import (
	"crypto/rand"
	"github.com/umeshdhaked/athens/pkg/logger"
	"io"
	"sync"
)

var (
	once    sync.Once
	service *OtpSender
)

type OtpSender struct {
}

func NewOtpSender() {
	once.Do(func() {
		service = &OtpSender{}
	})
}

func GetOtpSender() *OtpSender {
	return service
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
	logger.GetLogger().WithField("otp", otp).Info("otp sent")
	//send otp here
	return nil
}
