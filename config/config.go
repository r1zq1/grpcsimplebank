package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource    string `mapstructure:"DB_SOURCE"`
	GRPCAddress string `mapstructure:"GRPC_ADDRESS"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("warning: %v, fallback to env only", err)
	}
	err = viper.Unmarshal(&config)
	return config, err
}
