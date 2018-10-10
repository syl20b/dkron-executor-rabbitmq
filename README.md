# dkron-executor-rabbitmq

RabbitMQ Executor for [dkron](https://dkron.io) that publishes message in RabbitMQ queue.
This implementation is inspired from [bringg/dkron-executor-rabbitmq](https://github.com/bringg/dkron-executor-rabbitmq)

We have rewrite mainly because we need a different way to provide RabbitMQ configuration. 
Besides we had errors using `dep` to manage the dependencies, therefore we are using `vgo`.

## Usage

The application is available using the docker image `solocal/dkron-executor-rabbitmq`.

In addition of the Dkron configuration, you need to provide the url of the 
RabbitMQ. 

### Configuration file

To provide the configuration to the RabbitMQ executor, you can create a file 
`dkron-executor-rabbitmq.json` next to the Dkron configuration file. 

```
{
    "rabbit_connection_url": "amqp://guest:guest@localhost:5672/"
}
```

or 

```
{
    "rabbit_url": "localhost:5672"
    "rabbit_user": "guest"
    "rabbit_password": "guest"
}
```

### Environment variables

To provide the configuration to the RabbitMQ executor, you can also use 
environment variables :

- `DKRON_EXECUTOR_RABBITMQ_RABBIT_HOST` 

or 

- `DKRON_EXECUTOR_RABBITMQ_RABBIT_URL` 
- `DKRON_EXECUTOR_RABBITMQ_RABBIT_USER` 
- `DKRON_EXECUTOR_RABBITMQ_RABBIT_PASSWORD` 

### Create a job using the RabbitMQ executor

To create a job that publishes `a payload message` to the exchange `my-exchange` 
using the routing key `a.routing.key`, you can run

```bash
curl -n -X POST http://localhost:8080/v1/jobs \
 -H "Content-Type: application/json" \
 -d '{
   "name": "rabbit-job",
   "schedule": "@every 1m",
   "shell": false,
   "executor": "rabbitmq",
   "executor_config": {
     "exchange": "my-exchange",
     "routingKey": "a.routing.key",
     "payload": "a payload message"  
   },
   "disabled": false
}'
```



## Development

To run tests

```bash
./run-tests.sh
```

You can run test directly using `go test` but you will need to start DKron and RabbitMQ before.
To do this you can use `docker-compose run --rm start-dependencies` and `docker-compose run --rm start-services`
