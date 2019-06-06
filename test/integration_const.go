package integrationtest

import (
	"time"
	"logrusfluenthook/hook"
	"github.com/sirupsen/logrus"
)

var (
	contextTimeout = 20 * time.Second
	dockerOpsTimeout = 15 * time.Second
	loggerFieldMap = logrus.FieldMap{
		hook.FieldKeyData: "@data",
		logrus.FieldKeyTime:  "@timestamp",
		logrus.FieldKeyFile:  "@file",
		logrus.FieldKeyFunc:  "@func",
		logrus.FieldKeyLevel: "@level",
		logrus.FieldKeyMsg:   "@msg",
	}
	testTimestamp = time.Now().Format(time.RFC3339)
)

const (
	baseTag = "integration_testing"
	fluentConfigRelativePath = "./fluent_stdout.conf"
	fluentImageUrl = "docker.io/fluent/fluentd:v1.5-1"
	fluentContainerName = "logrus_fluent_hook_integration_test_instance"
	fluentRemotePort = "24224"
	localhost = "127.0.0.1"
)
