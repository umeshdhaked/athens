package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/config"
)

var (
	once sync.Once
	db   *Db
)

type Db struct {
	Client *dynamodb.Client
}

func NewDb() {
	once.Do(func() {
		// Create DynamoDB client
		cfg, err := awsConfig.LoadDefaultConfig(context.Background(),
			awsConfig.WithRegion(config.GetConfig().Aws.Db.Region),
			awsConfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: config.GetConfig().Aws.Db.EndPoint}, nil
				})))

		if err != nil {
			fmt.Println("Error loading AWS config:", err)
			return
		}

		client := dynamodb.NewFromConfig(cfg)

		db = &Db{
			Client: client,
		}
	})
}

func GetDb() *Db {
	return db
}
