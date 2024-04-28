package register

import (
	"context"
	"errors"
	"fmt"
	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/services/otp"
	"github.com/fastbiztech/hastinapura/internal/services/subscription"
	"log"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/pkg/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/jwt"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	once    sync.Once
	service *RegistrationService
)

type RegistrationService struct {
	otpService *otp.OtpService
	cryp       *crypto.Crypto
	userRepo   *repo.UserRepo
	sub        *subscription.SubscriptionService
}

func NewRegistrationService(userRepo *repo.UserRepo, otpService *otp.OtpService, cryp *crypto.Crypto, sub *subscription.SubscriptionService) {
	once.Do(func() {
		service = &RegistrationService{
			userRepo:   userRepo,
			otpService: otpService,
			cryp:       cryp,
			sub:        sub,
		}
	})
}

func GetRegistrationService() *RegistrationService {
	return service
}

func (s *RegistrationService) SendOtp(ctx *gin.Context, user dtos.RegisterUserRequest) error {
	err := s.otpService.SendOtp(ctx, user.MobileNumber)
	// add OTP send logic here.
	log.Println("otp sent for user ", user.MobileNumber)
	// save hashed otp after sending otp....
	return err
}

func (s *RegistrationService) UpdateUserRoleToAdmin(ctx *gin.Context, user dtos.RegisterUserRequest) (*dtos.LoginSuccessResponse, error) {
	log.Println("Received update user role to admin: " + user.MobileNumber)

	var usr *models.User
	var er error
	if usr, er = s.userRepo.GetUserFromMobile(ctx, user.MobileNumber); er != nil {
		return nil, er
	}

	usrObj := &models.User{ID: usr.ID, Role: "admin", Mobile: usr.Mobile}
	if er := s.userRepo.UpdateUser(ctx, usrObj); er != nil {
		return nil, er
	}

	return &dtos.LoginSuccessResponse{MobileNumber: user.MobileNumber}, nil
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

	usrObj := &models.User{ID: uuid.New().String(), Mobile: user.MobileNumber, Hashed_password: s.cryp.HashString(user.Password)}
	if er := s.userRepo.UpdateUser(ctx, usrObj); er != nil {
		return nil, er
	}

	// Add default subscriptions to user
	ctx.Set(constants.JwtTokenRole, "admin")
	ctx.Set(constants.JwtTokenUserID, "system")
	ctx.Set(constants.JwtTokenMobile, "1234567890")
	err := s.sub.AddDefaultSubscriptionToUser(ctx, &dtos.UserDefaultSubscriptionRequest{UserMobile: user.MobileNumber})
	if err != nil {
		return nil, err
	}

	token, err := jwt.CreateToken(usrObj.ID, usrObj.Mobile, usrObj.Role)
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
			token, err := jwt.CreateToken(usr.ID, usr.Mobile, usr.Role)
			if err != nil {
				return nil, err
			}

			return &dtos.LoginSuccessResponse{MobileNumber: usr.Mobile, LoginToken: token}, nil
		} else {
			return nil, errors.New("password did not match")
		}
	}
}

func (s *RegistrationService) GetUser(ctx context.Context, mobile string) (*models.User, error) {

	usr, err := s.userRepo.GetUserFromMobile(ctx, mobile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return usr, nil
}
