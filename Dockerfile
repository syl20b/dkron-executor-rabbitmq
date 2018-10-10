# build dkron-executor-rabbitmq
FROM golang:1.11.1 AS builder

RUN mkdir -p /app
WORKDIR /app

ENV GO111MODULE=on
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go install

# final image
FROM dkron/dkron:v0.11.0
LABEL maintainer="Leadformance <dev@leadformance.com>"

COPY --from=builder /go/bin/dkron-executor-rabbitmq /opt/local/dkron/dkron-executor-rabbitmq
