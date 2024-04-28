package controllers

import (
	"encoding/json"
	"errors"
	"github.com/fastbiztech/hastinapura/internal/services/payments"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go/utils"
	"io"
	"log"
	"net/http"
)

func HandlePaymentCreateOrder(ctx *gin.Context) {
	var orderReq dtos.PaymentOrderRequest
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var rzp *payments.PaymentService = payments.GetPaymentService()
	orderResp, err := rzp.CreateOrder(ctx, &orderReq)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, orderResp)
}

// Depricated, don't use HandleUpdatePaymentOrder, use HandlePaymentOrderWebhook
func HandleUpdatePaymentOrder(ctx *gin.Context) {
	var orderReq dtos.UpdatePaymentOrderRequest
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var rzp *payments.PaymentService = payments.GetPaymentService()
	orderResp, err := rzp.UpdatePaymentOrder(ctx, &orderReq)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, orderResp)
}

func HandlePaymentOrderWebhook(ctx *gin.Context) {
	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var orderReq dtos.PaymentWebhookRequest
	if err := json.Unmarshal(jsonData, &orderReq); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// validate webhook signature (request is coming from payments server only)
	sign := ctx.Request.Header["X-Razorpay-Signature"][0]
	isValid := utils.VerifyWebhookSignature(string(jsonData), sign, "23e12f50-3ee6-41b8-bcdb-fd123dfd28cb")
	if !isValid {
		log.Println("unable to validate signature with secret")
		ctx.String(http.StatusInternalServerError, errors.New("invalid signature").Error())
		return
	}

	rzp := payments.GetPaymentService()
	err = rzp.PaymentOrderWebhook(ctx, &orderReq)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, dtos.PaymentWebhookResponse{Status: "ok"})
}

func HandlePaymentStatus(ctx *gin.Context) {
	var req dtos.PaymentStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	rzp := payments.GetPaymentService()
	resp, err := rzp.GetPaymentStatus(ctx, req.OrderId)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandlePaymentsHistory(ctx *gin.Context) {
	var paymentHistoryReq dtos.PaymentHistoryRequest
	if err := ctx.ShouldBindJSON(&paymentHistoryReq); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	rzp := payments.GetPaymentService()
	payHistoryResp, err := rzp.GetPaymentsHistory(ctx, &paymentHistoryReq)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, payHistoryResp)
}
