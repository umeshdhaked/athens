package aws

import (
	"log"

	"github.com/FastBizTech/hastinapura/pkg/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ConfigureAwsSdkSession(config *models.ApplicationConfig) *session.Session {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AwsSecretConfig.Region),
		Credentials: credentials.NewStaticCredentials(config.AwsSecretConfig.AccessKeyID, config.AwsSecretConfig.SecretAccessKey, ""),
	})
	session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}))

	if err != nil {
		log.Fatal(err)
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		log.Fatal(err)
	}

	return sess
}
