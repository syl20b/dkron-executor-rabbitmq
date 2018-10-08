package main

import (
	"github.com/hashicorp/go-plugin"
	dkplugin "github.com/victorcoder/dkron/plugin"
)

func main() {
	config := readConfig()

	executor, err := createRabbitMQExecutor(config)
	if err != nil {
		log.WithError(err).Error("Failed to create rabbit mq executor")
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: dkplugin.Handshake,
		Plugins: map[string]plugin.Plugin{
			"executor": &dkplugin.ExecutorPlugin{Executor: executor},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
