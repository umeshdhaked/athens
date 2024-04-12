package crypto

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

var (
	once    sync.Once
	service *Crypto
)

type Crypto struct {
}

func NewCrypto() {
	once.Do(func() {
		service = &Crypto{}
	})
}

func GetCrypto() *Crypto {
	return service
}

func (o *Crypto) HashString(otp string) string {
	h := sha256.New()
	h.Write([]byte(string(otp)))

	bs := h.Sum(nil)
	s := fmt.Sprintf("%x\n", bs)
	return s
}
