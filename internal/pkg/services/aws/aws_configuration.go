package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	internalConfig "github.com/fastbiztech/hastinapura/internal/config"
	"log"
	"os"
)

func ConfigureAwsSdkConfig(internalConfig *internalConfig.Config) aws.Config {
	var cfg aws.Config
	var err error
	if os.Getenv("env") == "local" {
		cfg, err = config.LoadDefaultConfig(
			context.Background(),
			config.WithRegion(internalConfig.Aws.Db.Region),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: internalConfig.Aws.Db.EndPoint}, nil
			})),
			//config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			//	return aws.Credentials{
			//		AccessKeyID:     internalConfig.Aws.Db.KeyID,
			//		SecretAccessKey: internalConfig.Aws.Db.AccessKey,
			//		SessionToken:    "TOKEN",
			//	}, nil
			//})),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background())
	}

	if err != nil {
		log.Fatal("Error loading AWS config:", err)
	}
	return cfg
}
