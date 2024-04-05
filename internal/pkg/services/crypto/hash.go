package crypto

import (
	"crypto/sha256"
	"fmt"
)

type Crypto struct {
}

func NewCrypto() *Crypto {
	return &Crypto{}
}

func (o *Crypto) HashString(otp string) string {
	h := sha256.New()
	h.Write([]byte(string(otp)))

	bs := h.Sum(nil)
	s := fmt.Sprintf("%x\n", bs)
	return s
}
