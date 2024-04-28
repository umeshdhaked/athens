package payments

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/pkg/rzp"
	"github.com/fastbiztech/hastinapura/internal/services/subscription"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go/utils"
	"log"
	"sync"
)

var once sync.Once
var paymentService *PaymentService

type PaymentService struct {
	rzpService          *rzp.RzpService
	paymentRepo         *repo.Payments
	subscriptionService *subscription.SubscriptionService
}

func NewPaymentService(rzpService *rzp.RzpService, paymentRepo *repo.Payments, subscriptionService *subscription.SubscriptionService) {
	once.Do(func() {
		paymentService = &PaymentService{
			rzpService:          rzpService,
			paymentRepo:         paymentRepo,
			subscriptionService: subscriptionService,
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
	item, err := attributevalue.MarshalMap(body)
	if err != nil {
		return nil, err
	}
	err = r.paymentRepo.CreateItem(ctx, models.TableRzpOrders, item)
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
}

// UpdatePaymentOrder Deprecated
func (r *PaymentService) UpdatePaymentOrder(ctx *gin.Context, orderReq *dtos.UpdatePaymentOrderRequest) (*dtos.UpdatePaymentResponse, error) {
	//update payment table with data from UI
	paymentResponse, err := attributevalue.MarshalMap(orderReq)
	if err != nil {
		return nil, err
	}
	err = r.paymentRepo.CreateItem(ctx, models.TableFBTPayment, paymentResponse)
	if err != nil {
		return nil, err
	}

	// verify payment signature
	params := map[string]interface{}{
		"razorpay_order_id":   orderReq.RazorpayOrderId,
		"razorpay_payment_id": orderReq.RazorpayPaymentId,
	}
	isVerified := utils.VerifyPaymentSignature(params, orderReq.RazorpaySignature, models.TestPaymentSecret)
	if !isVerified {
		return nil, errors.New("signature is not correct")
	}

	//check here if order_id exists for user in RzpOrders table in created state
	usrId, _ := ctx.Get(constants.JwtTokenUserID)
	rzpOrder, err := r.paymentRepo.GetOrderFromId(ctx, orderReq.RazorpayOrderId)
	if err != nil {
		return nil, err
	}
	orderUser := rzpOrder.Notes.UserID
	if rzpOrder.Status != "created" {
		return nil, errors.New("order is not in created state")
	}
	if orderUser != usrId {
		return nil, errors.New("order is created by different user")
	}

	// get latest order from razorpay
	orderBody, err := r.rzpService.FetchOrder(orderReq.RazorpayOrderId)
	if err != nil {
		return nil, err
	}
	orderItem, err := attributevalue.MarshalMap(orderBody)
	if err != nil {
		return nil, err
	}
	err = r.paymentRepo.CreateItem(ctx, models.TableRzpOrders, orderItem)
	if err != nil {
		return nil, err
	}
	orderStatus, _ := orderBody["status"] // created, attempted, paid

	// get latest payment from razorpay
	paymentBody, er := r.rzpService.FetchPayment(orderReq.RazorpayPaymentId)
	if er != nil {
		return nil, er
	}
	paymentItem, err := attributevalue.MarshalMap(paymentBody)
	if err != nil {
		return nil, err
	}
	err = r.paymentRepo.CreateItem(ctx, models.TableRzpPayments, paymentItem)
	if err != nil {
		return nil, err
	}
	paymentStatus, _ := paymentBody["status"] // created authorized captured refunded failed

	amount, _ := paymentBody["amount"]
	mobile, _ := ctx.Get(constants.JwtTokenMobile)

	if orderStatus == "paid" && paymentStatus == "captured" {
		// payment is success, add credits to user.
		err := r.subscriptionService.AddCreditToUser(ctx, &dtos.AddCreditsRequest{
			UserMobile:     mobile.(string),
			InitialCredit:  amount.(float64) / 100,
			PaymentOrderId: rzpOrder.ID,
		})
		if err != nil {
			return nil, err
		}
	}

	// paid, captured is successful. or should if amount deducted will be refunded.
	return &dtos.UpdatePaymentResponse{OrderStatus: orderStatus.(string), PaymentStatus: paymentStatus.(string)}, nil
}

func (r *PaymentService) PaymentOrderWebhook(ctx *gin.Context, orderReq *dtos.PaymentWebhookRequest) error {
	orderBody := orderReq.Payload["order"].Entity

	rzpOrder, err := r.paymentRepo.GetOrderFromId(ctx, orderBody["id"].(string))
	if err != nil {
		return err
	}
	if rzpOrder.Status == "paid" {
		log.Println("idempotent request for order id")
		return nil
	}

	orderItem, err := attributevalue.MarshalMap(orderBody)
	if err != nil {
		return err
	}
	err = r.paymentRepo.CreateItem(ctx, models.TableRzpOrders, orderItem)
	if err != nil {
		return err
	}

	//// get latest payment from razorpay
	paymentBody := orderReq.Payload["payment"].Entity
	paymentItem, err := attributevalue.MarshalMap(paymentBody)
	if err != nil {
		return err
	}

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
			PaymentOrderId: rzpOrder.ID,
		})
		if err != nil {
			return err
		}
	}

	// update our order table only if credits are added/given to user for payment
	err = r.paymentRepo.CreateItem(ctx, models.TableRzpPayments, paymentItem)
	if err != nil {
		return err
	}

	return nil
}

func (r *PaymentService) GetPaymentStatus(ctx *gin.Context, orderId string) (string, error) {
	rzpOrder, err := r.paymentRepo.GetOrderFromId(ctx, orderId)
	if err != nil {
		return "", err
	}
	return rzpOrder.Status, nil
}

func (r *PaymentService) GetPaymentsHistory(ctx *gin.Context, req *dtos.PaymentHistoryRequest) (*dtos.PaymentHistoryResponse, error) {
	rzpOrder, lastEvaluatedKey, count, err := r.paymentRepo.GetOrderList(ctx, req.Limit, req.LastEvaluatedKey, req.Status)
	if err != nil {
		return nil, err
	}

	resp := &dtos.PaymentHistoryResponse{}
	resp.Count = count
	resp.LastEvaluatedKey = lastEvaluatedKey

	for _, r := range rzpOrder {
		resp.OrderList = append(resp.OrderList, struct {
			OrderId   string
			Amount    int
			CreatedAt int
		}{OrderId: r.ID, Amount: r.Amount, CreatedAt: r.CreatedAt})
	}

	return resp, nil
}

// get order history api mock status like [completed, failed] - get from order collection
// generate invoice ok payment, save data in invoice DB, create download invoice api
