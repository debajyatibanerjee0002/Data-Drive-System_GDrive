package config

import "github.com/spf13/viper"

func ViperConfig() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}
