package hook

import (
	"time"

	"github.com/sirupsen/logrus"
)

// DefaultFluentPort defines the default fluentd hostname if no
// host is specified in the configuration
const DefaultFluentHost = "0.0.0.0"

// DefaultFluentPort defines the default port to use if no
// port is specified in the configuration
const DefaultFluentPort uint16 = 24224

// DefaultFluentTagPrefix defines the default string used by fluentd
// as the base tag
const DefaultFluentTagPrefix = "app"

// DefaultLogTag is for the tag applied to all logs going to fluentd
// after the base tag
const DefaultLogTag = "log"

// DefaultTimestampFormat is for use when no timestamp format is defined
const DefaultTimestampFormat = time.RFC3339

// FieldKeyData is an extension of the logrus Fields for the data key
// so we can avoid having a DataKey property like logrus is doing
const FieldKeyData = "data"

// FieldMap is used for setting the correct properties on the logs
// being sent to fluentd
var DefaultFieldMap = map[string]string{
	logrus.FieldKeyFile:      "@file",
	logrus.FieldKeyFunc:      "@func",
	logrus.FieldKeyLevel:     "@level",
	logrus.FieldKeyMsg:   "@msg",
	logrus.FieldKeyTime: "@timestamp",
	FieldKeyData:      "@data",
}

// LevelMap is here so we can use string(s) instead of logrus.Level
// to define the levels at which to log
var LevelMap = map[string]interface{}{
	"trace": logrus.TraceLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
	"panic": logrus.PanicLevel,
}
