package register

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"

	"github.com/fastbiztech/hastinapura/api/services/otp"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
	"github.com/fastbiztech/hastinapura/internal/pkg/repositories"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/jwt"
	"github.com/fastbiztech/hastinapura/pkg/models/dtos"
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

func (s *RegistrationService) SendOtp(ctx *gin.Context, user dtos.RegisterUserRequest) error {
	err := s.otpService.SendOtp(ctx, user.MobileNumber)
	// add OTP send logic here.
	log.Println("otp sent for user ", user.MobileNumber)
	// save hashed otp after sending otp....
	return err
}

func (s *RegistrationService) RegisterUser(ctx *gin.Context, user dtos.RegisterUserRequest) (*dtos.LoginSuccessResponse, error) {
	log.Println("Received register user request for user ", user)
	if err := s.otpService.VerifyOtp(ctx, user.MobileNumber, user.Otp); err != nil {
		return nil, err
	}

	if usr, er := s.userRepo.GetUserFromMobile(ctx, user.MobileNumber); er != nil {
		return nil, er
	} else if usr != nil {
		return nil, errors.New("user already exists")
	}

	usrObj := &dbo.User{Id: uuid.New().String(), Mobile: user.MobileNumber, Hashed_password: s.cryp.HashString(user.Password)}
	if er := s.userRepo.CreateUser(ctx, usrObj); er != nil {
		return nil, er
	}

	token, err := jwt.CreateToken(usrObj.Id, usrObj.Mobile, usrObj.Role)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginSuccessResponse{MobileNumber: user.MobileNumber, LoginToken: token}, nil
}

func (s *RegistrationService) LoginUser(ctx *gin.Context, user dtos.RegisterUserRequest) (*dtos.LoginSuccessResponse, error) {
	mobile := user.MobileNumber
	password := user.Password

	usr, err := s.userRepo.GetUserFromMobile(ctx, mobile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		if s.cryp.HashString(password) == usr.Hashed_password {
			token, err := jwt.CreateToken(usr.Id, usr.Mobile, usr.Role)
			if err != nil {
				return nil, err
			}

			return &dtos.LoginSuccessResponse{MobileNumber: usr.Mobile, LoginToken: token}, nil
		} else {
			return nil, errors.New("password did not match")
		}
	}
}
