package rzp

import (
	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/razorpay/razorpay-go"
	"sync"
)

var once sync.Once
var rzpService *RzpService

type RzpService struct {
	client *razorpay.Client
}

func NewRzpService() {
	once.Do(func() {
		rzpService = &RzpService{
			client: razorpay.NewClient(models.TestPaymentKey, models.TestPaymentSecret),
		}
	})
}

func GetRzpService() *RzpService {
	return rzpService
}

func (r *RzpService) CreateOrder(ctx *gin.Context, amount int64) (map[string]interface{}, error) {
	usrId, _ := ctx.Get(constants.JwtTokenUserID)
	mobile, _ := ctx.Get(constants.JwtTokenMobile)
	data := map[string]interface{}{
		"amount":          amount,
		"currency":        models.CurrencyINR,
		"receipt":         uuid.New().String(),
		"partial_payment": false,
		"notes": map[string]interface{}{
			"userId": usrId,
			"mobile": mobile,
		},
	}
	return r.client.Order.Create(data, nil)
}

func (r *RzpService) FetchOrder(razorpayOrderId string) (map[string]interface{}, error) {
	return r.client.Order.Fetch(razorpayOrderId, nil, nil)
}

func (r *RzpService) FetchOrderPayment(orderId string) (map[string]interface{}, error) {
	return r.client.Order.Payments(orderId, nil, nil)
}

func (r *RzpService) FetchPayment(razorpayPaymentId string) (map[string]interface{}, error) {
	return r.client.Payment.Fetch(razorpayPaymentId, nil, nil)
}

func (r *RzpService) CreateRefund(paymentMap map[string]interface{}) (map[string]interface{}, error) {
	return r.client.Refund.Create(paymentMap, nil)
}
