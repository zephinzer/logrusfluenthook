package hook

import (
	"time"

	"github.com/sirupsen/logrus"
)

// FieldMap is used for setting the correct properties on the logs
// being sent to fluentd
var DefaultFieldMap = map[string]string{
	"timestamp": "@timestamp",
	"message":   "@msg",
	"data":      "@data",
	"caller":    "@caller",
	"level":     "@level",
	"file":      "@file",
	"func":      "@func",
}

// DefaultLogTag is for the tag applied to all logs going to fluentd
// after the base tag
const DefaultLogTag = "log"

// DefaultTimeFormat is for use when no timestamp format is defined
const DefaultTimeFormat = time.RFC3339

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

// var LogrusLevelMap = map[logrusLevel]string {
// 	logrus.TraceLevel: "trace",
// 	logrus.
// }
