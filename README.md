# dkron-executor-rabbitmq

RabbitMQ Executor for [dkron](https://dkron.io) that publishes message in RabbitMQ queue.
This implementation is inspired from [bringg/dkron-executor-rabbitmq](https://github.com/bringg/dkron-executor-rabbitmq)

We have rewrite mainly because we need a different way to provide RabbitMQ configuration. 
Besides we had errors using `dep` to manage the dependencies, therefore we are using `vgo`.