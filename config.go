package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	connectionUrl string
}

func readConfig() *Config {
	viper.SetDefault("rabbit_connection_url", "amqp://guest:guest@localhost:5672/")

	viper.SetEnvPrefix("dkron_executor_rabbitmq")
	viper.BindEnv("rabbit_connection_url")
	viper.BindEnv("rabbit_url")
	viper.BindEnv("rabbit_user")
	viper.BindEnv("rabbit_password")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		log.WithError(err).Info("No valid config found: Applying default values and Environment variable.")
	}

	return createConfig()
}

func createConfig() *Config {
	connectionUrl := getConnectionUrlFromEnvVar()
	if viper.IsSet("rabbit_url") {
		connectionUrl = buildConnectionUrlFromUrl()
	}

	config := &Config{
		connectionUrl: connectionUrl,
	}

	return config
}

func buildConnectionUrlFromUrl() string {
	rabbitUrl := viper.GetString("rabbit_url")

	rabbitUser := viper.GetString("rabbit_user")
	rabbitPassword := viper.GetString("rabbit_password")

	userPassword := ""
	if rabbitUser != "" && rabbitPassword != "" {
		userPassword = fmt.Sprintf("%s:%s@", rabbitUser, rabbitPassword)
	}

	connectionUrl := fmt.Sprintf("amqp://%s%s", userPassword, rabbitUrl)
	return connectionUrl
}

func getConnectionUrlFromEnvVar() string {
	return viper.GetString("rabbit_connection_url")
}
