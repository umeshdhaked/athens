package aws

import (
	"log"
	"os"

	"github.com/FastBizTech/hastinapura/pkg/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ConfigureAwsSdkSession(config *models.ApplicationConfig) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
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
