package config

import (
	"fmt"
	"log"
	"strings"
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

type BaseCronConfig struct {
	Enable        bool `mapstructure:"enable"`
	StartTime     int  `mapstructure:"start_time"`
	ExecutionTime int  `mapstructure:"execution_time"`
}

type CronsConfig struct {
	CronsConfigS3Contacts    BaseCronConfig `mapstructure:"s3_contacts"`
	CronsConfigPaymentRefund BaseCronConfig `mapstructure:"payment_refunds"`
}

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	Aws   AwsConfig   `mapstructure:"aws"`
	Api   ApiConfig   `mapstructure:"api"`
	Crons CronsConfig `mapstructure:"crons"`
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
	})
}

func GetConfig() *Config {
	return config
}
