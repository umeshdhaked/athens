package aws

import (
	"context"
	"fmt"
	"mime/multipart"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fastbiztech/hastinapura/internal/config"
)

const (

	// Specify the bucket name
	BucketContactUpload = "myfirstbucket"
)

type S3Client struct {
	s3Client *s3.Client
}

var (
	once     sync.Once
	s3Client *S3Client
)

func InitialiseS3Client() {
	once.Do(
		func() {
			// Create DynamoDB client
			cfg, err := awsConfig.LoadDefaultConfig(context.Background(),
				awsConfig.WithRegion(config.GetConfig().Aws.S3.Region),
				awsConfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
					func(service, region string, options ...interface{}) (aws.Endpoint, error) {
						return aws.Endpoint{URL: config.GetConfig().Aws.S3.EndPoint}, nil
					})))
			if err != nil {
				fmt.Println("Error loading AWS config:", err)
				return
			}

			// Create S3 client
			s3Client = &S3Client{}

			// Create client with Path style option instead of virtual host style.
			s3Client.s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.UsePathStyle = true
			})
		},
	)
}

func (c *S3Client) Upload(file multipart.File, bucketName, objectKey string) error {
	// Upload the file to S3
	_, err := c.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("error uploading file to S3: %v", err)
	}

	return nil
}

func GetS3Client() *S3Client {
	return s3Client
}
