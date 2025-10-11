package kyc

import (
	"github.com/umeshdhaked/athens/internal/constants"
	"github.com/umeshdhaked/athens/internal/models"
	"github.com/umeshdhaked/athens/internal/pkg/aws"
	"github.com/umeshdhaked/athens/internal/pkg/repo"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/umeshdhaked/athens/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"strconv"
	"sync"
	"time"
)

var (
	once       sync.Once
	kycService *KycService
)

type KycService struct {
	baseRepo repo.IRepository
}

func NewKycService() {
	once.Do(func() {
		kycService = &KycService{baseRepo: repo.GetRepository()}
	})
}

func GetKycService() *KycService {
	return kycService
}

func (k *KycService) UploadKycDocuments(ctx *gin.Context, docFile multipart.File, photo multipart.File, kycRequest *dtos.UploadKycRequest) (*dtos.UploadKycResponse, error) {
	docKey, photoKey := uuid.New().String(), uuid.New().String()
	// Upload file to S3
	if err := aws.GetS3Client().Upload(docFile, aws.BucketKycUpload, docKey); err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}
	if err := aws.GetS3Client().Upload(photo, aws.BucketKycUpload, photoKey); err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	kycDBO := models.Kyc{
		UserId:       strconv.FormatInt(ctx.GetInt64(constants.JwtTokenUserID), 10),
		UserName:     ctx.GetString(constants.JwtUserName),
		Mobile:       ctx.GetString(constants.JwtTokenMobile),
		DocumentType: kycRequest.DocType,
		KycDocLink:   docKey,
		PhotoLink:    photoKey,
		IsVerified:   "false",
		BaseModel:    models.BaseModel{CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
	}

	if err := k.baseRepo.Create(ctx, &kycDBO); err != nil {
		return nil, err
	}

	return &dtos.UploadKycResponse{Status: "success"}, nil
}

func (k *KycService) GetPendingKycVerifications(ctx *gin.Context, kycRequest *dtos.PendingKycsRequest) ([]*dtos.PendingKycsResponse, error) {

	kycDBO := []*models.Kyc{}
	err := k.baseRepo.FindMultiplePagination(ctx, &kycDBO, map[string]interface{}{
		"is_verified": "false",
	}, kycRequest.Pagination)

	if err != nil {
		return nil, err
	}

	resp := []*dtos.PendingKycsResponse{}
	for _, k := range kycDBO {
		// Photo Read from S3
		photoReader, err := aws.GetS3Client().Fetch(ctx, aws.BucketKycUpload, k.PhotoLink)
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("error fetching s3 file")
			return nil, err
		}
		photo, err := io.ReadAll(photoReader)
		if err != nil {
			return nil, err
		}
		// Doc read from S3
		docReader, err := aws.GetS3Client().Fetch(ctx, aws.BucketKycUpload, k.KycDocLink)
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("error fetching s3 file")
			return nil, err
		}
		doc, err := io.ReadAll(docReader)
		if err != nil {
			return nil, err
		}
		resp = append(resp, &dtos.PendingKycsResponse{
			UserName:    k.UserName,
			Mobile:      k.Mobile,
			DocType:     k.DocumentType,
			Photo:       photo,
			KycDocument: doc,
			KycId:       k.ID,
		})
	}
	// Test the response files Photo, KycDocument using https://base64.guru/converter/decode/file (Base64 to File decoder)

	return resp, nil
}

func (k *KycService) UpdateKycStatus(ctx *gin.Context, kycRequest *dtos.KycStatusUpdateRequest) (*dtos.KycStatusUpdateResponse, error) {
	kycDBO := &models.Kyc{}
	if err := k.baseRepo.FindByID(ctx, kycRequest.KycId, kycDBO); err != nil {
		return nil, err
	}

	userId, _ := strconv.Atoi(kycDBO.UserId)
	if kycRequest.Status == "approved" {
		user := &models.User{ID: int64(userId), KycDone: "true"}
		if err := k.baseRepo.Update(ctx, user); err != nil {
			return nil, err
		}
	}
	kycDBO.IsVerified = kycRequest.Status
	kycDBO.Comment = kycRequest.Comment
	if err := k.baseRepo.Update(ctx, kycDBO); err != nil {
		return nil, err
	}

	kycResp := &dtos.KycStatusUpdateResponse{
		KycId:  kycDBO.ID,
		Status: kycDBO.IsVerified,
	}

	return kycResp, nil
}

func (k *KycService) GetUserKycStatus(ctx *gin.Context, kycRequest *dtos.KycStatusRequest) (*dtos.KycStatusResponse, error) {
	kycDBO := &models.Kyc{}
	if err := k.baseRepo.FindByID(ctx, kycRequest.UserId, kycDBO); err != nil {
		return nil, err
	}

	kycResp := &dtos.KycStatusResponse{
		KycId:   kycDBO.ID,
		Status:  kycDBO.IsVerified,
		Comment: kycDBO.Comment,
	}

	return kycResp, nil
}
