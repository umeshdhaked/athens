package models

type AwsSecretConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

type ApplicationConfig struct {
	Port            string
	AwsSecretConfig AwsSecretConfig
}
