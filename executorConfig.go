package main

import (
	"fmt"
	"strconv"
)

type ExecutorConfig struct {
	exchange   string
	routingKey string
	mandatory  bool
	immediate  bool

	payload     string
	contentType string
}

func NewExecutorConfig(job string, configMap map[string]string) *ExecutorConfig {
	config := new(ExecutorConfig)

	config.routingKey = configMap["routing_key"]
	config.exchange = readFacultativeString(job, configMap, "exchange", "")
	config.mandatory = readFacultativeBool(job, configMap, "mandatory", false)
	config.immediate = readFacultativeBool(job, configMap, "immediate", false)

	config.payload = configMap["payload"]
	config.contentType = readFacultativeString(job, configMap, "contentType", "application/json")

	return config
}

func readFacultativeString(job string, configMap map[string]string, key string, defaultValue string) string {
	if value, ok := configMap[key]; ok {
		return value
	}

	return defaultValue
}

func readFacultativeBool(job string, configMap map[string]string, key string, defaultValue bool) bool {
	if value, ok := configMap[key]; ok {
		param, err := strconv.ParseBool(value)

		if err != nil {
			log.WithField("jobName", job).
				WithError(err).
				Error(fmt.Sprintf("Invalid value for '%s' field", key))

			return defaultValue
		}

		return param
	}

	return defaultValue
}
