package services

import (
	"log"

	"github.com/FastBizTech/hastinapura/pkg/models"
)

type RegistrationService struct {
}

func NewRegistrationService() *RegistrationService {
	return &RegistrationService{}
}

func (s *RegistrationService) SendOtp(user models.RegisterUserRequest) (bool, error) {
	// add OTP send logic here.
	log.Println("otp sent for user ", user.MobileNumber)
	return true, nil
}

func (s *RegistrationService) RegisterUser(user models.RegisterUserRequest) (*models.LoginSuccessResponse, error) {
	log.Println("Received register user request for user ", user)
	return &models.LoginSuccessResponse{user.MobileNumber, "some-login-token-after-signup"}, nil
}

func (s *RegistrationService) LoginUser(user models.RegisterUserRequest) (*models.LoginSuccessResponse, error) {
	mobile := user.MobileNumber
	password := user.Password
	log.Println("Login User  ", mobile+password)
	return &models.LoginSuccessResponse{user.MobileNumber, "some-login-token"}, nil
}
