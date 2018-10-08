package main

import (
	"github.com/streadway/amqp"
	"github.com/victorcoder/dkron/dkron"
)

type RabbitMQExecutor struct {
	connectionUrl string
	conn          *amqp.Connection
	ch            *amqp.Channel
}

func createRabbitMQExecutor(config *Config) (*RabbitMQExecutor, error) {
	executor := &RabbitMQExecutor{connectionUrl: config.connectionUrl}
	if err := executor.connect(); err != nil {
		return nil, err
	}

	return executor, nil
}

func (s *RabbitMQExecutor) connect() error {
	log.WithField("connectionUrl", s.connectionUrl).Info("Connecting to rabbit")

	conn, err := amqp.Dial(s.connectionUrl)
	if err != nil {
		log.WithField("connectionUrl", s.connectionUrl).Error("Fail to connect to rabbitmq")
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Error("Fail to open channel")
		return err
	}

	s.conn = conn
	s.ch = ch
	return nil
}

func (s *RabbitMQExecutor) publish(args *dkron.ExecuteRequest) error {
	executorConfig := NewExecutorConfig(args.JobName, args.Config)

	return s.ch.Publish(
		executorConfig.exchange,
		executorConfig.routingKey,
		executorConfig.mandatory,
		executorConfig.immediate,
		amqp.Publishing{
			ContentType: executorConfig.contentType,
			Body:        []byte(executorConfig.payload),
		})

}
func (s *RabbitMQExecutor) Execute(args *dkron.ExecuteRequest) ([]byte, error) {
	err := s.publish(args)
	if err == amqp.ErrClosed {
		log.WithField("jobName", args.JobName).Debug("Got closed error while trying to publish, trying to reconnect")

		if err = s.connect(); err != nil {
			log.WithField("jobName", args.JobName).WithError(err).Error("Failed to reconnect")
			return []byte(err.Error()), err
		}

		err = s.publish(args)
	}
	if err != nil {
		return []byte(err.Error()), err
	}

	return []byte("OK"), nil
}
