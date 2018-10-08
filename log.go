package main

import "github.com/sirupsen/logrus"

var log = logrus.NewEntry(logrus.New()).WithField("executor", "DKRON_EXECUTOR_RABBITMQ")
