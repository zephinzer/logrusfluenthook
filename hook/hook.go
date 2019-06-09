package hook

import (
	"fmt"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/sirupsen/logrus"
)

// New creates a new instance of the Hook for logrus to use;
// create it and add it to your logrus instance using:
//
//   hook, err := logrusfluenthook.New(&logrusfluenthook.Config{/* ... */})
//   if err != nil {
//     // handle the error
//   } else {
//     logrusInstance.AddHook(hook)
//   }
func New(config *Config) (*Hook, error) {
	levels := []logrus.Level{}

	config.assignDefaults()

	// create fluent logger instance
	fluentInstance, err := fluent.New(fluent.Config{
		FluentHost:         config.Host,
		FluentPort:         int(config.Port),
		TagPrefix:          config.BaseTag,
		SubSecondPrecision: true,
		RequestAck:         true,
	})
	if err != nil {
		if errorContainsString(err, "connection refused") {
			return nil, fmt.Errorf("fluentd does not seem to be available at %s:%v", config.Host, config.Port)
		}
		return nil, err
	}

	// create hook instance
	hook := Hook{
		config:         config,
		fluentInstance: fluentInstance,
		levels:         levels,
	}

	hook.SetLevels(config.Levels)

	return &hook, nil
}

// Hook is well, the hook; the logrus.Hook interface is implemented here
type Hook struct {
	config         *Config
	fluentInstance *fluent.Fluent
	levels         []logrus.Level
}

// getFieldName retrieves the customised label for the fluentd logs,
// if the requested :field is unassigned, it falls back to the default
// stored in DefaultFieldMap
func (hook *Hook) getFieldName(field string) string {
	if hook.config.FieldMap[field] != "" {
		return hook.config.FieldMap[field]
	}
	return DefaultFieldMap[field]
}

// getLogTag retrieves the tag to append to the base tag for fluentd
// to parse
func (hook Hook) getLogTag(entry *logrus.Entry) string {
	if entry.Data["tag"] != nil {
		if tag, ok := entry.Data["tag"].(string); ok {
			return tag
		}
	}
	return entry.Level.String()
}

// getTimeFormat returns the time format according to the configuration,
// if that's not available, returns the default time format
func (hook Hook) getTimeFormat() string {
	if hook.config.TimestampFormat != "" {
		return hook.config.TimestampFormat
	}
	return DefaultTimestampFormat
}

// getLogData retrieves the data from the provided :entry for use in
// the log sent to fluentd
func (hook *Hook) getLogData(entry *logrus.Entry) map[string]interface{} {
	logData := make(map[string]interface{})
	logData[hook.getFieldName(logrus.FieldKeyTime)] = entry.Time.Format(hook.getTimeFormat())
	logData[hook.getFieldName(logrus.FieldKeyMsg)] = entry.Message
	logData[hook.getFieldName(FieldKeyData)] = map[string]interface{}(entry.Data)
	logData[hook.getFieldName(logrus.FieldKeyLevel)] = entry.Level.String()
	if entry.HasCaller() {
		logData[hook.getFieldName(logrus.FieldKeyFile)] = fmt.Sprintf("%s:%v", entry.Caller.File, entry.Caller.Line)
		logData[hook.getFieldName(logrus.FieldKeyFunc)] = entry.Caller.Function
	}
	return logData
}

// Fire implements the logrus.Hook interface which is triggered when a log is fired
func (hook *Hook) Fire(entry *logrus.Entry) error {
	tag := hook.getLogTag(entry)
	data := hook.getLogData(entry)
	return hook.fluentInstance.Post(tag, data)
}

// Levels implements the logrus.Hook interface and returns the levels
// for which the Hook should be triggered
func (hook Hook) Levels() []logrus.Level {
	return hook.levels
}

// SetField sets the label of the provided :fieldName where :fieldName
// is one of "timestamp", "caller", "data", and "message"
func (hook *Hook) SetField(fieldName, fieldLabel string) {
	hook.config.FieldMap[fieldName] = fieldLabel
}

// SetLevels resets the levels to the provided :levels parameter
func (hook *Hook) SetLevels(levels []string) error {
	hook.levels = []logrus.Level{}
	if len(levels) == 0 {
		for _, logrusLevel := range LevelMap {
			hook.levels = append(hook.levels, logrusLevel.(logrus.Level))
		}
	} else {
		for _, stringLevel := range levels {
			if LevelMap[stringLevel] != nil {
				hook.levels = append(hook.levels, LevelMap[stringLevel].(logrus.Level))
			}
		}
	}
	return nil
}

// SetLogrusLevels is for if you wish to not use the convenience strings
func (hook *Hook) SetLogrusLevels(levels []logrus.Level) {
	hook.levels = levels
}
