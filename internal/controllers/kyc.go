package controllers

import (
	"github.com/umeshdhaked/athens/internal/services/kyc"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleUploadKycDocuments(c *gin.Context) {
	// Parse multipart form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request = dtos.UploadKycRequest{}
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kycDocument, _, err := c.Request.FormFile("KycDocument")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer kycDocument.Close()

	photo, _, err := c.Request.FormFile("Photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer photo.Close()

	var kycService *kyc.KycService = kyc.GetKycService()
	kycResp, err := kycService.UploadKycDocuments(c, kycDocument, photo, &request)
	if nil != err {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, kycResp)
}

func HandleGetPendingKycVerifications(ctx *gin.Context) {
	var kycReq dtos.PendingKycsRequest
	if err := ctx.ShouldBindJSON(&kycReq); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var kycService *kyc.KycService = kyc.GetKycService()
	kycResp, err := kycService.GetPendingKycVerifications(ctx, &kycReq)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, kycResp)
}

func HandleUpdateKycStatus(ctx *gin.Context) {
	var kycReq dtos.KycStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&kycReq); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var kycService *kyc.KycService = kyc.GetKycService()
	kycResp, err := kycService.UpdateKycStatus(ctx, &kycReq)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, kycResp)
}

func HandleGetKycStatus(ctx *gin.Context) {
	var kycReq dtos.KycStatusRequest
	if err := ctx.ShouldBindJSON(&kycReq); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var kycService *kyc.KycService = kyc.GetKycService()
	kycResp, err := kycService.GetUserKycStatus(ctx, &kycReq)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, kycResp)
}
