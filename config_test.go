package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setupTestCase(envvars map[string]string) func(t *testing.T) {
	for env, value := range envvars {
		os.Setenv(env, value)
	}

	return func(t *testing.T) {
		for env := range envvars {
			os.Unsetenv(env)
		}
	}
}

func TestReadConfig(t *testing.T) {
	cases :=
		[]struct {
			name     string
			envvars  map[string]string
			expected *Config
		}{
			{
				"withConnectionString",
				map[string]string{
					"DKRON_EXECUTOR_RABBITMQ_RABBIT_CONNECTION_URL": "amqp://localhost:5672",
				},
				&Config{
					connectionUrl: "amqp://localhost:5672",
				},
			},
			{
				"withUrlUserPassword",
				map[string]string{
					"DKRON_EXECUTOR_RABBITMQ_RABBIT_URL":      "localhost:5672",
					"DKRON_EXECUTOR_RABBITMQ_RABBIT_USER":     "myuser",
					"DKRON_EXECUTOR_RABBITMQ_RABBIT_PASSWORD": "mypassword",
				},
				&Config{
					connectionUrl: "amqp://myuser:mypassword@localhost:5672",
				},
			},
			{
				"withUrl",
				map[string]string{
					"DKRON_EXECUTOR_RABBITMQ_RABBIT_URL": "localhost:5672",
				},
				&Config{
					connectionUrl: "amqp://localhost:5672",
				},
			},
		}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := setupTestCase(tc.envvars)
			defer teardownSubTest(t)

			result := readConfig()
			assert.Equal(t, tc.expected, result)
		})
	}
}
