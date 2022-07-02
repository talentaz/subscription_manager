package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"Server"`
	Database struct {
		DBDriver   string `yaml:"driver"`
		DBUsername string `yaml:"DBUsername"`
		DBName     string `yaml:"DBName"`
		DBPassword string `yaml:"DBPassword"`
		DBSchema   string `yaml:"DBSchema"`
	} `yaml:"Database"`
	Stripe struct {
		StripeAPI      string `yaml:"StripeAPI"`
		CancelURL      string `yaml:"CancelURL"`
		SuccessURL     string `yaml:"SuccessURL"`
		EndpointSecret string `yaml:"EndpointSecret"`
	} `yaml:"Stripe"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("subscription_manager")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
