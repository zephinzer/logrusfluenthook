package integrationtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"
	"strconv"
	"testing"
	"time"

	"logrusfluenthook/hook"

	dockerContainer "github.com/docker/docker/api/types/container"
	dockerMount "github.com/docker/docker/api/types/mount"
	dockerClient "github.com/docker/docker/client"
	nat "github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/sirupsen/logrus"
)

// TestIntegration is the entrypoint for the test suite
func TestIntegration(t *testing.T) {
	suite.Run(t, &IntegrationTestSuite{})
}

// IntegrationTestSuite is the suite that contains tests for the actual
// integration of the hook with a fluentd service
type IntegrationTestSuite struct {
	docker            *dockerClient.Client
	fluentConfigPath string
	fluentLocalPort   int
	fluentContainer dockerContainer.ContainerCreateCreatedBody

	logger *logrus.Logger
	applicationLogs bytes.Buffer
	fluentLogs bytes.Buffer
	suite.Suite
}

// TestIntegration is the main test driver
func (s *IntegrationTestSuite) TestIntegration() {
	// send a log to the logger instance
	s.createSampleLog()
	latestFluentLog := s.retrieveLatestFluentLog()
	s.Contains(latestFluentLog, testTimestamp)

	// normalise container logs
	log.Println("normalising logs from remote fluentd and application...")
	expectedLogTag := fmt.Sprintf("%s.info", baseTag)
	startOfFluentData :=
		strings.Index(latestFluentLog, expectedLogTag) +
		len(expectedLogTag) +
		len(": ")
	fluentLogOfInterest := latestFluentLog[startOfFluentData:]
	appLogOfInterest := s.applicationLogs.String()

	// compare logs
	log.Println("comaparing remote log vs local log...")
	var unmarshalledFluentLog map[string]interface{}
	var unmarshalledAppLog map[string]interface{}
	if err := json.Unmarshal([]byte(fluentLogOfInterest), &unmarshalledFluentLog); err != nil {
		s.Error(err)
	} else if  err := json.Unmarshal([]byte(appLogOfInterest), &unmarshalledAppLog); err != nil {
		s.Error(err)
	}
	log.Printf("> application log line:\n\n%s\n\n", unmarshalledAppLog)
	log.Printf("> fluentd log line:\n\n%s\n\n", unmarshalledFluentLog)

	s.Equal(unmarshalledFluentLog, unmarshalledAppLog)

	log.Println("looks like all is well!")
}

func (s *IntegrationTestSuite) BeforeTest(suiteName, testName string) {
	log.Printf("# BeforeTest - %s.%s", suiteName, testName)

	log.Println("setting up a new logger instance...")
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetOutput(&s.applicationLogs)
	logger.SetFormatter(&logrus.JSONFormatter{
		DataKey: "@data",
		FieldMap: loggerFieldMap,
	})

	log.Println("setting up a new fluentd hook...")
	hookConfig := hook.Config{
		Host:    localhost,
		Port:    uint16(s.fluentLocalPort),
		BaseTag: baseTag,
		Levels:  []string{"trace", "debug", "info", "warn", "error", "panic"},
	}
	log.Printf("  - configuration: %v", hookConfig)
	logrushook, err := hook.New(&hookConfig)
	if err != nil {
		panic(err)
	}

	log.Println("  - attaching fluentd hook to logger instance...")
	logger.AddHook(logrushook)

	log.Println("assigning created logger to the global logger")
	s.logger = logger
}

func (s *IntegrationTestSuite) SetupTest() {
	log.Println("# SetupTest")
	var err error

	s.fluentConfigPath = path.Join(getCurrentFileDirectory(), fluentConfigRelativePath)
	log.Printf("fluent config path: %s\n", s.fluentConfigPath)

	log.Println("instantiating a docker client...")
	s.docker = createDockerClient()

	log.Printf("pulling fluent image from '%s'...", fluentImageUrl)
	pullImage(s.docker, fluentImageUrl)

	log.Printf("finding a free local port... ")
	s.fluentLocalPort = getFreePort()
	log.Printf("> found free port '%v'", s.fluentLocalPort)

	log.Printf("creating container '%s'...", fluentContainerName)
	s.fluentContainer, err = s.docker.ContainerCreate(
		context.Background(),
		&dockerContainer.Config{
			ExposedPorts: nat.PortSet{
				fluentRemotePort: struct{}{},
			},
			Image: fluentImageUrl,
			Tty:   false,
		},
		&dockerContainer.HostConfig{
			AutoRemove: true,
			Mounts: []dockerMount.Mount{
				{
					Type: dockerMount.TypeBind,
					Source: s.fluentConfigPath,
					Target: "/fluentd/etc/fluent.conf", // fluentd's default
				},
			},
			PortBindings: map[nat.Port][]nat.PortBinding{
				fluentRemotePort: {{HostIP: localhost, HostPort: strconv.Itoa(s.fluentLocalPort)}},
			},
		},
		nil,
		fluentContainerName,
	)
	if err != nil {
		panic(err)
	} else if len(s.fluentContainer.ID) == 0 {
		panic("container id should not have been blank, something went wrong")
	}
	log.Printf("> container created with id '%s'", s.fluentContainer.ID)

	log.Printf("starting container '%s'...", s.fluentContainer.ID)
	startContainer(s.docker, s.fluentContainer.ID)
	log.Println("> fluentd service has been created, waiting for 5 seconds for fluentd to initialise before proceeding...")

	<-time.After(5 * time.Second)
}

func (s *IntegrationTestSuite) TearDownTest() {
	log.Println("# TearDownTest")

	log.Printf("stopping container '%s'", s.fluentContainer.ID)
	stopContainer(s.docker, s.fluentContainer.ID)
	log.Println("> fluentd container has been stopped")
}

func (s *IntegrationTestSuite) retrieveLatestFluentLog() string {
	log.Println("# retrieveLatestFluentLog")

	log.Printf("retrieving container logs for container id '%s'...", s.fluentContainer.ID)
	fluentContainerLogs := getContainerLogs(s.docker, s.fluentContainer.ID)
	log.Printf("> retrieved the following logs from fluent:\n\n# BEGIN fluent container logs - - - -\n\n%s\n\n# ENDOF fluent container logs - - - -\n\n", fluentContainerLogs)
	
	linesOfLogs, lastLineOfLogs := getLastLineOfString(fluentContainerLogs)
	log.Printf("processed %v lines of logs", linesOfLogs)
	log.Printf("last retrieved line:\n\n# BEGIN last log fluentd received\n\n%s\n\n# ENDOF last log fluentd received\n\n", lastLineOfLogs)

	return lastLineOfLogs
}

func (s *IntegrationTestSuite) createSampleLog() {
	log.Println("# createSampleLog")
	logFields := logrus.Fields{
		"field_string":  "1",
		"field_integer": 1,
		"field_float":   1.01,
		"field_bool":    true,
		"field_time":    testTimestamp,
	}
	log.Printf("creating a sample line of log with the following data: %v", logFields)

	s.logger.WithFields(logFields).Info("hello world")
}
