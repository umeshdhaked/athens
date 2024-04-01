package aws

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fastbiztech/hastinapura/internal/config"
)

func ConfigureAwsSdkSession(config *config.Config) *session.Session {
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
