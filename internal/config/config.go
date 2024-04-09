package config

import (
	"fmt"
	"log"
	"strings"
	"os"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/spf13/viper"
)

var (
	config *Config
	once   sync.Once
)

type AwsDbConfig struct {
	EndPoint  string `mapstructure:"endpoint"`
	KeyID     string `mapstructure:"key_id"`
	AccessKey string `mapstructure:"access_key"`
	Region    string `mapstructure:"region"`
}

type AwsS3Config struct {
	EndPoint  string `mapstructure:"endpoint"`
	KeyID     string `mapstructure:"key_id"`
	AccessKey string `mapstructure:"access_key"`
	Region    string `mapstructure:"region"`
}

type AwsConfig struct {
	Db AwsDbConfig `mapstructure:"db"`
	S3 AwsS3Config `mapstructure:"s3"`
}

type AuthConfig struct {
	Key    string `mapstructure:"key"`
	Secret string `mapstructure:"secret"`
}

type EndpointConfig struct {
	Method  string     `mapstructure:"method"`
	BaseUrl string     `mapstructure:"base_url"`
	Path    string     `mapstructure:"path"`
	Auth    AuthConfig `mapstructure:"auth"`
}

type ApiConfig struct {
	SmsHeader  EndpointConfig `mapstructure:"sms_header"`
	InstantSms EndpointConfig `mapstructure:"instant_sms"`
}

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	App AppConfig `mapstructure:"app"`
	Aws AwsConfig `mapstructure:"aws"`
	Api ApiConfig `mapstructure:"api"`
}

func LoadConfig() {
	once.Do(func() {
		viper.AddConfigPath(utils.GetFilePath("config"))
		// Set the environment variable prefix
		viper.SetEnvPrefix("ENV_VAR")
		// Enable automatic environment variable binding
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// Load values from default file first.
		viper.SetConfigName("default")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal(err.Error())
		}

		// Load env specific values and merge.
		viper.SetConfigName(utils.GetEnv())

		err = viper.MergeInConfig() // This will merge env.json with default.json
		if err != nil {
			log.Fatalf("Error reading env config: %s", err)
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			fmt.Printf("unable to decode into config struct, %v", err)
			panic("failed to load config")
		}

		aws_id := os.Getenv("AWS_ACCESS_KEY_ID")
		aws_sec := os.Getenv("AWS_ACCESS_KEY_ID")
		aws_region := os.Getenv("AWS_REGION")
		if aws_id != "" {
			if config.Aws.Db.KeyID != "" {
				config.Aws.Db.KeyID = aws_id
			}
			if config.Aws.S3.KeyID != "" {
				config.Aws.S3.KeyID = aws_id
			}
		}
		if aws_sec != "" {
			if config.Aws.Db.AccessKey != "" {
				config.Aws.Db.AccessKey = aws_id
			}
			if config.Aws.S3.AccessKey != "" {
				config.Aws.S3.AccessKey = aws_id
			}
		}
		if aws_region != "" {
			if config.Aws.Db.Region != "" {
				config.Aws.Db.Region = aws_id
			}
			if config.Aws.S3.Region != "" {
				config.Aws.S3.Region = aws_id
			}
		}
	})
}

func GetConfig() *Config {
	return config
}
