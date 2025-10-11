package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/umeshdhaked/athens/internal/config"
	"github.com/umeshdhaked/athens/pkg/logger"
)

func dynomoDbInit() *dynamodb.Client {
	// Create DynamoDB client
	cfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(config.GetConfig().Db.Dynamo.Region),
		awsConfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: config.GetConfig().Db.Dynamo.EndPoint,
				}, nil
			}),
		),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.GetConfig().Db.Dynamo.KeyID,
				config.GetConfig().Db.Dynamo.AccessKey,
				"", // session token, if any
			),
		),
	)
	if err != nil {
		logger.GetLogger().WithField("message", err.Error()).Panic("failed to initialise dynomodb")
	}

	client := dynamodb.NewFromConfig(cfg)

	// ping check
	var apiErr *types.ResourceNotFoundException
	_, err = client.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String("Otp"), // Replace with your table name
	})
	if err != nil && !errors.As(err, &apiErr) {
		logger.GetLogger().WithField("message", err.Error()).Panic("ping failed")
	}

	logger.GetLogger().Info("dynamo db connection successful")

	return client
}
