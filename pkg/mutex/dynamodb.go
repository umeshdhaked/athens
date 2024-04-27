package mutex

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"cirello.io/dynamolock/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TODO: build releaseLocks functionality for expired locks (for any ungraceful shutdown)

const (
	TableLocks         = "locks"
	AcquireLockTimeout = 1 * time.Second

	S3ContactsFetchProcessingLeaseDuration   = 30 * time.Second
	S3ContactsFetchProcessingStartAtDuration = 0 * time.Second

	PaymentRefundProcessingLeaseDuration   = 30 * time.Second
	PaymentRefundProcessingStartAtDuration = 0 * time.Second
)

var (
	once                         sync.Once
	s3ProcessingMutex            *DynamoDBLockManager
	paymentRefundProcessingMutex *DynamoDBLockManager
)

// handlerFunc defines the interface for handler function
type HandlerFunc func() (interface{}, error)

// Lock represents a lock acquired from the lock manager.
type Lock struct {
	lockItem *dynamolock.Lock
}

// DynamoDBLockManager manages distributed locks in DynamoDB
type DynamoDBLockManager struct {
	TableName string
	Client    *dynamolock.Client
}

func Initialise() {
	once.Do(func() {
		var err error

		// s3 processing mutex instance
		s3ProcessingMutex, err = NewDynamoDBLockManager(S3ContactsFetchProcessingLeaseDuration, S3ContactsFetchProcessingStartAtDuration)
		paymentRefundProcessingMutex, err = NewDynamoDBLockManager(PaymentRefundProcessingLeaseDuration, PaymentRefundProcessingStartAtDuration)

		if err != nil {
			log.Fatal("failed initialising all mutexes")
		}
	})
}

func GetS3ProcessingMutexLockManager() *DynamoDBLockManager {
	return s3ProcessingMutex
}

func PaymentRefundProcessingMutexLockManager() *DynamoDBLockManager {
	return paymentRefundProcessingMutex
}

// NewDynamoDBLockManager creates a new DynamoDBLockManager instance
func NewDynamoDBLockManager(leaseDuration, heartbeatPeriod time.Duration) (*DynamoDBLockManager, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.Background(),
		awsConfig.WithRegion(config.GetConfig().Aws.Db.Region),
		awsConfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: config.GetConfig().Aws.Db.EndPoint}, nil
			})),
	)
	if err != nil {
		log.Fatal("Error loading AWS config:", err)
	}

	client, err := dynamolock.New(dynamodb.NewFromConfig(cfg),
		TableLocks,
		dynamolock.WithLeaseDuration(leaseDuration),
		dynamolock.WithHeartbeatPeriod(heartbeatPeriod),
		//dynamolock.WithLogger(log.New(os.Stderr, "", 0)), // TODO: update logger
	)
	if err != nil {
		return nil, err
	}

	return &DynamoDBLockManager{
		TableName: TableLocks,
		Client:    client,
	}, nil
}

func (lm *DynamoDBLockManager) AcquireAndRelease(ctx *gin.Context, key string, data []byte, handler HandlerFunc) (
	interface{}, error) {

	lockItem, err := lm.acquireLock(key, data)
	if err != nil {
		log.Println("failed to acquire lock: %w", err)
		return nil, err
	}

	// This is mandatory to release the mutex once the handler execution is done.
	// even when there is any runtime error/panic we have to release the mutex
	// so that it can be acquired by other process
	defer func() {
		err = lm.releaseLock(lockItem)
		if err != nil {
			log.Println("failed to release  lock: %w", err)
		}
	}()

	return handler()
}

// AcquireLock acquires a lock for the given key and data.
func (lm *DynamoDBLockManager) acquireLock(key string, data []byte) (Lock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), AcquireLockTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			// If the context times out, return an error indicating lock acquisition failure.
			return Lock{}, errors.New("failed to acquire lock: timeout exceeded")
		default:
			// Attempt to acquire the lock.
			lockItem, err := lm.Client.AcquireLock(key,
				dynamolock.WithData(data),
				dynamolock.ReplaceData(),
				dynamolock.FailIfLocked(),
				dynamolock.WithDeleteLockOnRelease(),
			)
			if err != nil {
				// If an error occurs during lock acquisition, retry after a delay.
				return Lock{}, err
			}
			// Lock acquired successfully, return the Lock object.
			return Lock{lockItem}, nil
		}
	}
}

// ReleaseLock releases the lock.
func (lm *DynamoDBLockManager) releaseLock(lock Lock) error {
	success, err := lm.Client.ReleaseLock(lock.lockItem)
	if !success {
		return errors.New("lock was lost before release")
	}
	return err
}

func ConnectCheck() {
	newLockManager, err := NewDynamoDBLockManager(20*time.Second, 1*time.Second)
	if err != nil {
		log.Fatal("failed to initialise mutex client")
	}

	/*
		test connection
		Create a dummy key and data for the lock.
	*/
	dummyKey := uuid.New().String()

	// Attempt to acquire a dummy lock.
	lockItem, err := newLockManager.Client.AcquireLock(dummyKey,
		dynamolock.ReplaceData(),
		dynamolock.FailIfLocked(),
		dynamolock.WithDeleteLockOnRelease(),
	)
	if err != nil {
		log.Fatal("failed to acquire dummy lock: %w", err)
	}

	// Release the dummy lock.
	_, err = newLockManager.Client.ReleaseLock(lockItem)
	if err != nil {
		log.Fatal("failed to release dummy lock: %w", err)
	}

	return
}
