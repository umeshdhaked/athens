package mutex

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

// HandlerFunc defines the interface for handler function
type HandlerFunc func() (interface{}, error)

var (
	once   sync.Once
	client *Mutex
)

type Mutex struct {
	r *redsync.Redsync
}

func Initialise() {
	once.Do(func() {
		NewRedsync()
	})
}

func GetClient() *Mutex {
	return client
}

func NewRedsync() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GetConfig().Mutex.Host, config.GetConfig().Mutex.Port),
		Password: config.GetConfig().Mutex.Password,
		DB:       0,
	})

	pool := goredis.NewPool(redisClient)

	client = &Mutex{
		r: redsync.New(pool),
	}

	// Check ping
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.GetLogger().WithField("message", err.Error()).Panic("failed to initialise redis mutex")
	}
}

func (m *Mutex) AcquireAndRelease(ctx *gin.Context, key string, ttl time.Duration, handler HandlerFunc) (interface{}, error) {
	mutex := m.r.NewMutex(key, redsync.WithExpiry(ttl))

	// Acquire lock
	if err := mutex.LockContext(ctx); err != nil {
		return nil, err
	}

	// Release lock when done
	defer func() {
		_, err := mutex.UnlockContext(ctx)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return
		}
		logger.GetLogger().Info("mutex released successfully for key: " + key)
	}()

	// Perform protected operations here
	return handler()
}
