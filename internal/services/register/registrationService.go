package register

import (
	"errors"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/services/otp"
	"github.com/fastbiztech/hastinapura/internal/services/subscription"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	gormLogger "gorm.io/gorm/logger"

	"github.com/fastbiztech/hastinapura/internal/pkg/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/jwt"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

var (
	once    sync.Once
	service *RegistrationService
)

type RegistrationService struct {
	otpService *otp.OtpService
	cryp       *crypto.Crypto
	baseRepo   repo.IRepository
	userRepo   repo.IUserRepo
	sub        *subscription.SubscriptionService
}

func NewRegistrationService(otpService *otp.OtpService, cryp *crypto.Crypto, sub *subscription.SubscriptionService) {
	once.Do(func() {
		service = &RegistrationService{
			baseRepo:   repo.GetRepository(),
			userRepo:   repo.GetUserRepo(),
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
	logger.GetLogger().WithField("mobile", user.MobileNumber).Info("otp sent for user")
	// save hashed otp after sending otp....
	return err
}

func (s *RegistrationService) UpdateUserRoleToAdmin(ctx *gin.Context, request dtos.RegisterUserRequest) (*dtos.LoginSuccessResponse, error) {
	logger.GetLogger().WithField("mobile", request.MobileNumber).Info("Received update user role to admin: ")

	var (
		user *models.User = &models.User{}
		err  error
	)

	// get User
	err = s.baseRepo.Find(ctx, user, map[string]interface{}{
		models.ColumnUserMobile: request.MobileNumber,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	user.Role = "admin"
	err = s.baseRepo.Update(ctx, user)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return &dtos.LoginSuccessResponse{MobileNumber: request.MobileNumber}, nil
}

func (s *RegistrationService) RegisterUser(ctx *gin.Context, request dtos.RegisterUserRequest) (*dtos.LoginSuccessResponse, error) {
	logger.GetLogger().WithField("request", request).Info("Received register user request for user ")
	if err := s.otpService.VerifyOtp(ctx, request.MobileNumber, request.Otp); err != nil {
		return nil, err
	}

	// get User
	var user models.User
	err := s.baseRepo.Find(ctx, &user, map[string]interface{}{
		models.ColumnUserMobile: request.MobileNumber,
	})

	if err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound) {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	if user.ID != 0 {
		logger.GetLogger().Error("user already exists")
		return nil, errors.New("user already exists")
	}

	user = models.User{
		Mobile:          request.MobileNumber,
		Hashed_password: s.cryp.HashString(request.Password),
	}
	if err = s.baseRepo.Create(ctx, &user); err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Add default subscriptions to user
	ctx.Set(constants.JwtTokenRole, "admin")
	ctx.Set(constants.JwtTokenUserID, int64(-1))
	ctx.Set(constants.JwtTokenMobile, "1234567890")
	err = s.sub.AddDefaultSubscriptionToUser(ctx, &dtos.UserDefaultSubscriptionRequest{UserMobile: request.MobileNumber})
	if err != nil {
		return nil, err
	}

	token, err := jwt.CreateToken(user.ID, user.Mobile, user.Role)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginSuccessResponse{MobileNumber: request.MobileNumber, LoginToken: token}, nil
}

func (s *RegistrationService) LoginUser(ctx *gin.Context, request dtos.RegisterUserRequest) (*dtos.LoginSuccessResponse, error) {

	// get User
	var user models.User
	err := s.baseRepo.Find(ctx, &user, map[string]interface{}{
		models.ColumnUserMobile: request.MobileNumber,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	if s.cryp.HashString(request.Password) == user.Hashed_password {
		token, err := jwt.CreateToken(user.ID, user.Mobile, user.Role)
		if err != nil {
			return nil, err
		}

		return &dtos.LoginSuccessResponse{MobileNumber: user.Mobile, LoginToken: token}, nil
	} else {
		return nil, errors.New("password did not match")
	}

}

func (s *RegistrationService) GetUser(ctx *gin.Context, mobile string) (*models.User, error) {

	// get User
	var user models.User
	err := s.baseRepo.Find(ctx, &user, map[string]interface{}{
		models.ColumnUserMobile: mobile,
	})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return &user, nil
}
