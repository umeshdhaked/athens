package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
)

func ConfigureAwsSdkSession() *session.Session {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatal(err)
	}
	// t, err := sess.Config.Credentials.Get()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// print(t.AccessKeyID)

	return sess
}
