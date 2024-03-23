package models

import (
	"encoding/json"
	"log"
	"os"

	"context"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AwsSecretManager struct {
	SecretName    string
	TestingSecret string
}

type ApplicationConfig struct {
	Port             string
	AwsSecretManager AwsSecretManager
}

func (a *ApplicationConfig) UpdateSecretsInConfig(secretCache *secretcache.Cache) {

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(a.AwsSecretManager.SecretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	var secretsJsonString string = *result.SecretString

	mp := make(map[string]string)
	json.Unmarshal([]byte(secretsJsonString), &mp)

	setSecretToConfig(a, mp)
}

func setSecretToConfig(config *ApplicationConfig, mp map[string]string) {
	config.AwsSecretManager.TestingSecret = mp["AwsSecretConfig.TestingSecret"]
}
