package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_driver            string        `mapstructure:"DB_DRIVER"`
	DB_source            string        `mapstructure:"DB_SOURCE"`
	HttpServerAddr       string        `mapstructure:"HTTP_SADDR"`
	GrpcServerAddr       string        `mapstructure:"GRPC_SADDR"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccesTokenDuration   time.Duration `mapstructure:"ACCES_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func InitConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
