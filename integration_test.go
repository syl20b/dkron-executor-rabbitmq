package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"text/template"
	"time"
)

type Job struct {
	Name       string
	Schedule   string
	Exchange   string
	RoutingKey string
	Payload    string
}

const Lifetime = 5 * time.Second
const Exchange = "test-exchange"
const ExchangeType = "direct"
const Queue = "test-queue"
const RoutingKey = "my.routing.key"

var job = &Job{
	Name:       "test-job",
	Schedule:   "@every 1s",
	Exchange:   Exchange,
	RoutingKey: RoutingKey,
	Payload:    "trololo",
}

func rabbitUrl() string {
	rabbitUrl := os.Getenv("RABBIT_URL")

	if rabbitUrl == "" {
		return "amqp://localhost:5672"
	}

	return rabbitUrl
}

func dkronUrl() string {
	dkronUrl := os.Getenv("DKRON_URL")

	if dkronUrl == "" {
		return "http://localhost:8080"
	}

	return dkronUrl
}

func TestRabbitExecutor(t *testing.T) {
	handler := func(consumer *Consumer, deliveries <-chan amqp.Delivery, done chan error) {
		for d := range deliveries {
			log.WithField("message", string(d.Body)).Info("Message received")

			assert.Equal(t, job.Payload, string(d.Body))
			d.Ack(false)

			log.Debug("shutting down")

			deleteJob(t, job)
			if err := consumer.Shutdown(); err != nil {
				log.WithError(err).Fatal("error during shutdown")
			}
		}
		log.Debug("handle: deliveries channel closed")
		done <- nil
	}

	_, err := NewConsumer(rabbitUrl(), Exchange, ExchangeType, Queue, RoutingKey, handler)
	if err != nil {
		log.WithError(err).Fatal("Can't start the consumer")
	}

	createJob(t, job)

	time.Sleep(Lifetime)
}

func createJob(t *testing.T, job *Job) {
	jobTmpl := `{   
"name": "{{.Name}}",   
"schedule": "{{.Schedule}}",   
"shell": false,   
"executor": "rabbitmq",   
"executor_config": {     
	"exchange": "{{.Exchange}}",     
	"queue_name": "{{.RoutingKey}}",
	"payload": "{{.Payload}}"   
},   
"disabled": false 
}`

	jobTemplate := template.Must(template.New("job").Parse(jobTmpl))
	jsonStr := &bytes.Buffer{}
	if err := jobTemplate.Execute(jsonStr, job); err != nil {
		panic(err)
	}

	_, err := http.Post(fmt.Sprintf("%s/v1/jobs", dkronUrl()), "encoding/json", jsonStr)
	if err != nil {
		t.Fatal(err)
	}
}

func deleteJob(t *testing.T, job *Job) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/jobs/%s", dkronUrl(), job.Name), nil)
	if err != nil {
		t.Fatal(err)
	}

	client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}
