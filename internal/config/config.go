package config

import (
	"fmt"

	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/spf13/viper"
)

var (
	config *Config
)

type AwsDbConfig struct {
	EndPoint  string
	KeyID     string
	AccessKey string
	Region    string
}

type AwsConfig struct {
	Db AwsDbConfig
}

type AppConfig struct {
	Port string
}

type Config struct {
	App AppConfig
	Aws AwsConfig
}

func LoadConfig() {
	viper.AddConfigPath(utils.GetFilePath("api/config"))
	viper.SetConfigName(utils.GetEnv())

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		panic("failed to load config")
	}
}

func GetConfig() *Config {
	return config
}
