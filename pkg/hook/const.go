package hook

import "github.com/sirupsen/logrus"

// DefaultLogTag is for the tag applied to all logs going to fluentd
// after the base tag
const DefaultLogTag = "log"

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

// FieldMap is used for setting the correct properties on the logs
// being sent to fluentd
var FieldMap = map[string]string{
	"timestamp": "@timestamp",
	"message":   "@msg",
	"data":      "@data",
	"caller":    "@caller",
}
