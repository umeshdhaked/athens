package register

import (
	"errors"
	"fmt"
	"log"

	"github.com/fastbiztech/hastinapura/api/services/otp"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/requests"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/responses"
	"github.com/fastbiztech/hastinapura/internal/pkg/repositories"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/jwt"
	"github.com/google/uuid"
)

type RegistrationService struct {
	otpService *otp.OtpService
	cryp       *crypto.Crypto
	userRepo   *repositories.UserRepo
}

func NewRegistrationService(userRepo *repositories.UserRepo, otpService *otp.OtpService, cryp *crypto.Crypto) *RegistrationService {
	return &RegistrationService{userRepo: userRepo, otpService: otpService, cryp: cryp}
}

func (s *RegistrationService) SendOtp(user requests.RegisterUserRequest) error {
	err := s.otpService.SendOtp(user.MobileNumber)
	// add OTP send logic here.
	log.Println("otp sent for user ", user.MobileNumber)
	// save hashed otp after sending otp....
	return err
}

func (s *RegistrationService) RegisterUser(user requests.RegisterUserRequest) (*responses.LoginSuccessResponse, error) {
	log.Println("Received register user request for user ", user)
	if err := s.otpService.VerifyOtp(user.MobileNumber, user.Otp); err != nil {
		return nil, err
	}

	if usr, er := s.userRepo.GetUserFromMobile(user.MobileNumber); er != nil {
		return nil, er
	} else if usr != nil {
		return nil, errors.New("user already exists")
	}

	usrObj := &dbo.User{Id: uuid.New().String(), Mobile: user.MobileNumber, Hashed_password: s.cryp.HashString(user.Password)}
	if er := s.userRepo.CreateUser(usrObj); er != nil {
		return nil, er
	}

	token, err := jwt.CreateToken(usrObj.Id, usrObj.Mobile, usrObj.Role)
	if err != nil {
		return nil, err
	}

	return &responses.LoginSuccessResponse{MobileNumber: user.MobileNumber, LoginToken: token}, nil
}

func (s *RegistrationService) LoginUser(user requests.RegisterUserRequest) (*responses.LoginSuccessResponse, error) {

	mobile := user.MobileNumber
	password := user.Password

	usr, err := s.userRepo.GetUserFromMobile(mobile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		if s.cryp.HashString(password) == usr.Hashed_password {
			token, err := jwt.CreateToken(usr.Id, usr.Mobile, usr.Role)
			if err != nil {
				return nil, err
			}

			return &responses.LoginSuccessResponse{MobileNumber: usr.Mobile, LoginToken: token}, nil
		} else {
			return nil, errors.New("password did not match")
		}

	}
}
