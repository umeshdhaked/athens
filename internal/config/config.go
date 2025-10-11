package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/fastbiztech/hastinapura/pkg/logger"

	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/spf13/viper"
)

const (
	// MysqlConnectionDSNFormat : DSN for connecting mysql
	MysqlConnectionDSNFormat = "%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local"
)

var (
	config *Config
	once   sync.Once
)

type AwsDbConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
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

type DbConfig struct {
	Mysql  MysqlConfig `mapstructure:"mysql"`
	Dynamo AwsDbConfig `mapstructure:"dynamo"`
}

type MysqlConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Dialect  string `mapstructure:"dialect"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Protocol string `mapstructure:"protocol"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type MutexConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type BaseCronConfig struct {
	Enable        bool `mapstructure:"enable"`
	StartTime     int  `mapstructure:"start_time"`
	ExecutionTime int  `mapstructure:"execution_time"`
}

type CronsConfig struct {
	CronsConfigCampaign      BaseCronConfig `mapstructure:"campaign"`
	CronsConfigS3Contacts    BaseCronConfig `mapstructure:"s3_contacts"`
	CronsConfigPaymentRefund BaseCronConfig `mapstructure:"payment_refunds"`
}

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	Db    DbConfig    `mapstructure:"db"`
	Mutex MutexConfig `mapstructure:"mutex"`
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
			logger.GetLogger().Panic(err.Error())
		}

		// Load env specific values and merge.
		viper.SetConfigName(utils.GetEnv())

		err = viper.MergeInConfig() // This will merge env.json with default.json
		if err != nil {
			logger.GetLogger().WithField("error", err).Panic("Error reading env config")
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("unable to decode into config struct")
			logger.GetLogger().Panic("failed to load config")
		}

		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err == nil {
			logger.GetLogger().Info(string(configBytes))
		}
	})
}

func GetConfig() *Config {
	return config
}

func (c DbConfig) URL() string {
	// charset=utf8: uses utf8 character set data format
	// parseTime=true: changes the output type of DATE and DATETIME values to time.Time instead of []byte / strings
	// loc=Local: Sets the location for time.Time values (when using parseTime=true). "Local" sets the system's location
	switch c.Mysql.Dialect {
	case "mysql":
		return fmt.Sprintf(
			MysqlConnectionDSNFormat,
			c.Mysql.UserName,
			c.Mysql.Password,
			c.Mysql.Protocol,
			c.Mysql.Host,
			c.Mysql.Port,
			c.Mysql.Name)
	default:
		return ""
	}
}
