package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBDriver   string `yaml:"driver"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	DBUsername string `yaml:"DBUsername"`
	DBName     string `yaml:"DBName"`
	DBPassword string `yaml:"DBPassword"`
	DBSchema   string `yaml:"DBSchema"`
	StripeAPI  string `yaml:"StripeAPI"`
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
