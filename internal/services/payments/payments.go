package payments

import (
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"sync"
	"time"

	"github.com/fastbiztech/hastinapura/internal/models"

	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/pkg/rzp"
	"github.com/fastbiztech/hastinapura/internal/services/subscription"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

var once sync.Once
var paymentService *PaymentService

type PaymentService struct {
	rzpService          *rzp.RzpService
	baseRepo            repo.IRepository
	paymentRepo         repo.IPaymentsRepo
	invoiceRepo         repo.IInvoiceRepo
	subscriptionService *subscription.SubscriptionService
}

func NewPaymentService(rzpService *rzp.RzpService, paymentRepo repo.IPaymentsRepo, invoiceRepo repo.IInvoiceRepo, subscriptionService *subscription.SubscriptionService) {
	once.Do(func() {
		paymentService = &PaymentService{
			rzpService:          rzpService,
			baseRepo:            repo.GetRepository(),
			paymentRepo:         paymentRepo,
			subscriptionService: subscriptionService,
			invoiceRepo:         invoiceRepo,
		}
	})
}

func GetPaymentService() *PaymentService {
	return paymentService
}

func (r *PaymentService) CreateOrder(ctx *gin.Context, orderReq *dtos.PaymentOrderRequest) (*dtos.PaymentOrderResponse, error) {
	body, err := r.rzpService.CreateOrder(ctx, orderReq.Amount)
	if err != nil {
		return nil, err
	}
	rzpOrderDBO := &models.Payments{}
	rzpOrderDBO.PopulateFromMap(body)

	err = r.baseRepo.Create(ctx, rzpOrderDBO)
	if err != nil {
		return nil, err
	}
	//log.Println(body)

	orderId, _ := body["id"]
	amount, _ := body["amount"]

	return &dtos.PaymentOrderResponse{
		Key:        models.TestPaymentKey,
		RzpOrderId: orderId.(string),
		Amount:     amount.(float64),
		Currency:   models.CurrencyINR,
		OrgName:    models.PaymentOrgName,
	}, nil

	return nil, nil
}

// UpdatePaymentOrder Deprecated
func (r *PaymentService) UpdatePaymentOrder(ctx *gin.Context, orderReq *dtos.UpdatePaymentOrderRequest) (*dtos.UpdatePaymentResponse, error) {
	////update payment table with data from UI
	//paymentResponse, err := attributevalue.MarshalMap(orderReq)
	//if err != nil {
	//	return nil, err
	//}
	//err = r.paymentRepo.CreateItem(ctx, models.TableFBTPayment, paymentResponse)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// verify payment signature
	//params := map[string]interface{}{
	//	"razorpay_order_id":   orderReq.RazorpayOrderId,
	//	"razorpay_payment_id": orderReq.RazorpayPaymentId,
	//}
	//isVerified := utils.VerifyPaymentSignature(params, orderReq.RazorpaySignature, models.TestPaymentSecret)
	//if !isVerified {
	//	return nil, errors.New("signature is not correct")
	//}
	//
	////check here if order_id exists for user in RzpOrders table in created state
	//usrId, _ := ctx.Get(constants.JwtTokenUserID)
	//rzpOrder, err := r.paymentRepo.GetOrderFromId(ctx, orderReq.RazorpayOrderId)
	//if err != nil {
	//	return nil, err
	//}
	//orderUser := rzpOrder.Notes.UserID
	//if rzpOrder.Status != "created" {
	//	return nil, errors.New("order is not in created state")
	//}
	//if orderUser != usrId {
	//	return nil, errors.New("order is created by different user")
	//}
	//
	//// get latest order from razorpay
	//orderBody, err := r.rzpService.FetchOrder(orderReq.RazorpayOrderId)
	//if err != nil {
	//	return nil, err
	//}
	//orderItem, err := attributevalue.MarshalMap(orderBody)
	//if err != nil {
	//	return nil, err
	//}
	//err = r.paymentRepo.CreateItem(ctx, models.TableRzpOrders, orderItem)
	//if err != nil {
	//	return nil, err
	//}
	//orderStatus, _ := orderBody["status"] // created, attempted, paid
	//
	//// get latest payment from razorpay
	//paymentBody, er := r.rzpService.FetchPayment(orderReq.RazorpayPaymentId)
	//if er != nil {
	//	return nil, er
	//}
	//paymentItem, err := attributevalue.MarshalMap(paymentBody)
	//if err != nil {
	//	return nil, err
	//}
	//err = r.paymentRepo.CreateItem(ctx, models.TableRzpPayments, paymentItem)
	//if err != nil {
	//	return nil, err
	//}
	//paymentStatus, _ := paymentBody["status"] // created authorized captured refunded failed
	//
	//amount, _ := paymentBody["amount"]
	//mobile, _ := ctx.Get(constants.JwtTokenMobile)
	//
	//if orderStatus == "paid" && paymentStatus == "captured" {
	//	// payment is success, add credits to user.
	//	err := r.subscriptionService.AddCreditToUser(ctx, &dtos.AddCreditsRequest{
	//		UserMobile:     mobile.(string),
	//		InitialCredit:  amount.(float64) / 100,
	//		PaymentOrderId: rzpOrder.ID,
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//
	//// paid, captured is successful. or should if amount deducted will be refunded.
	//return &dtos.UpdatePaymentResponse{OrderStatus: orderStatus.(string), PaymentStatus: paymentStatus.(string)}, nil

	return nil, nil
}

func (r *PaymentService) PaymentOrderWebhook(ctx *gin.Context, orderReq *dtos.PaymentWebhookRequest) error {
	orderBody := orderReq.Payload["order"].Entity
	existingOrder := &models.Payments{}
	condition := make(map[string]interface{})
	condition["order_id"] = orderBody["id"]
	err := r.baseRepo.Find(ctx, existingOrder, condition)
	if err != nil {
		return err
	}
	if existingOrder.Status == "paid" {
		logger.GetLogger().Error("idempotent request for order id")
		return nil
	}

	rzpOrderBody := &models.Payments{}
	rzpOrderBody.PopulateFromMap(orderBody)
	rzpOrderBody.Id = existingOrder.Id

	//// get latest payment from razorpay
	paymentBody := orderReq.Payload["payment"].Entity
	// paymentItem, err := attributevalue.MarshalMap(paymentBody)
	// if err != nil {
	// 	return err
	// }

	orderStatus, _ := orderBody["status"]     // created, attempted, paid
	paymentStatus, _ := paymentBody["status"] // created authorized captured refunded failed
	amount, _ := paymentBody["amount"]
	notes, _ := orderBody["notes"]
	mobile, _ := notes.(map[string]interface{})["mobile"]

	if orderStatus == "paid" && paymentStatus == "captured" {
		// payment is success, add credits to user.
		err := r.subscriptionService.AddCreditToUser(ctx, &dtos.AddCreditsRequest{
			UserMobile:     mobile.(string),
			InitialCredit:  amount.(float64) / 100,
			PaymentOrderId: rzpOrderBody.OrderId,
		})
		if err != nil {
			return err
		}
	}

	err = r.baseRepo.Update(ctx, rzpOrderBody)
	if err != nil {
		return err
	}
	// err = r.paymentRepo.CreateItem(ctx, models.TableRzpPayments, paymentItem)
	// if err != nil {
	// 	return err
	// }

	//Get Empty Invoice
	invoice := &models.Invoice{}
	invoice.OrderId = rzpOrderBody.OrderId
	invoice.Status = "CREATED"
	invoice.UserId = rzpOrderBody.UserId
	invoice.Receipt = rzpOrderBody.Receipt
	invoice.BaseModel = models.BaseModel{CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()}
	err = r.baseRepo.Create(ctx, invoice)
	if err != nil {
		return err
	}

	return nil
}

func (r *PaymentService) GetPaymentStatus(ctx *gin.Context, orderId string) (string, error) {
	rzpOrder := &models.Payments{}
	err := r.baseRepo.Find(ctx, rzpOrder, map[string]interface{}{
		models.SQLColumnInvoiceOrderId: orderId,
	})
	if err != nil {
		return "", err
	}
	return rzpOrder.Status, nil
}

func (r *PaymentService) GetPaymentsHistory(ctx *gin.Context, req *dtos.PaymentHistoryRequest) (*dtos.PaymentHistoryResponse, error) {
	rzpOrders := []*models.Payments{}
	err := r.baseRepo.FindMultiplePagination(ctx, &rzpOrders, map[string]interface{}{
		"user_id": req.UserId,
		"status":  req.Status,
	}, req.Pagination)

	if err != nil {
		return nil, err
	}

	resp := &dtos.PaymentHistoryResponse{}
	resp.Orders = rzpOrders

	return resp, nil
}

// get order history api mock status like [completed, failed] - get from order collection
// generate invoice ok payment, save data in invoice DB, create download invoice api

//func (r *PaymentService) GenerateInvoice(ctx *gin.Context) {
//	getEmpty, er := r.invoiceRepo.GetEmptyInvoice(ctx)
//	if er != nil {
//		return
//	}
//	fmt.Print(currentInvoiceNumber)
//}
