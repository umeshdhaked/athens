package payments

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/pkg/rzp"
	"github.com/fastbiztech/hastinapura/pkg/cron"
	"github.com/fastbiztech/hastinapura/pkg/mutex"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"time"
)

const (
	S3ContactsFetchProcessingBatchSize = 5

	// mutex lock keys
	MutexKeyPendingPaymentFetchProcessing = "mutex-key-payment-fetch-processing"
)

var oncePaymentCron sync.Once
var paymentCronService *PaymentCronService

type PaymentCronService struct {
	ctx         *gin.Context
	rzpService  *rzp.RzpService
	paymentRepo *repo.Payments
}

func NewPaymentCronService(rzpService *rzp.RzpService, paymentRepo *repo.Payments) {
	oncePaymentCron.Do(func() {
		newCtx, _ := gin.CreateTestContext(nil)
		paymentCronService = &PaymentCronService{
			ctx:         newCtx,
			rzpService:  rzpService,
			paymentRepo: paymentRepo,
		}
	})
}

func GetPaymentCronService() *PaymentCronService {
	return paymentCronService
}

func InitiateRefundForStuckOrdersCron(p *PaymentCronService) {
	//if os.Getenv(constants.WorkerCronArg) != constants.WorkerCronArgPaymentRefund {
	//	return
	//}

	if !config.GetConfig().Crons.CronsConfigPaymentRefund.Enable {
		return
	}

	job := (&cron.Scheduler{}).NewScheduler()

	job.Initialize(
		time.Duration(config.GetConfig().Crons.CronsConfigPaymentRefund.ExecutionTime)*time.Second,
		time.Duration(config.GetConfig().Crons.CronsConfigPaymentRefund.StartTime)*time.Second,
		p,
	)
}

func (r *PaymentCronService) JobExecutor() {

	_, err := mutex.PaymentRefundProcessingMutexLockManager().
		AcquireAndRelease(r.ctx,
			MutexKeyPendingPaymentFetchProcessing,
			[]byte("Dummy Data"),
			func() (interface{}, error) {
				createdAtBefore := time.Now().Unix() - 4*60*60 // check for order created 4 hour ago

				createdOrder, err := r.paymentRepo.GetCreatedOrders(context.Background(), createdAtBefore)
				if err != nil {
					return nil, err
				}
				attemptedOrder, err := r.paymentRepo.GetAttemptedOrders(context.Background(), createdAtBefore)
				if err != nil {
					return nil, err
				}
				orders := append(createdOrder, attemptedOrder...)
				for _, order := range orders {
					orderBody, err := r.rzpService.FetchOrder(order.ID)
					if err != nil {
						return nil, err
					}

					payments, err := r.rzpService.FetchOrderPayment(order.ID)
					if err != nil {
						return nil, err
					}
					for _, payment := range payments["items"].([]interface{}) {
						paymentId := payment.(map[string]interface{})["id"].(string)
						paymentStatus := payment.(map[string]interface{})["status"].(string)

						paymentMap := payment.(map[string]interface{})
						paymentMap["payment_id"] = paymentId
						if paymentStatus == "captured" {
							rfnd, err := r.rzpService.CreateRefund(paymentMap) // maybe create a refund table
							if err != nil {
								return nil, err
							}
							rfndItem, err := attributevalue.MarshalMap(rfnd)
							if err != nil {
								return nil, err
							}
							err = r.paymentRepo.CreateItem(context.Background(), models.TableRzpRefunds, rfndItem)
							if err != nil {
								return nil, err
							}
						}
						// get payment again
						paymentBody, err := r.rzpService.FetchPayment(paymentId)
						if err != nil {
							return nil, err
						}
						paymentItem, err := attributevalue.MarshalMap(paymentBody)
						if err != nil {
							return nil, err
						}
						err = r.paymentRepo.CreateItem(context.Background(), models.TableRzpPayments, paymentItem)
						if err != nil {
							return nil, err
						}
					}

					// mark order as cancel/refunded
					orderBody["status"] = "refunded/canceled"
					orderItem, err := attributevalue.MarshalMap(orderBody)
					if err != nil {
						return nil, err
					}
					err = r.paymentRepo.CreateItem(context.Background(), models.TableRzpOrders, orderItem)
					if err != nil {
						return nil, err
					}
				}
				return nil, nil
			})

	if nil != err {
		log.Println(err)
	}
}
