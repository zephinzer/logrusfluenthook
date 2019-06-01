package main

import (
	logrusfluenthook "logrusfluenthook/pkg/hook"

	"github.com/sirupsen/logrus"
)

func main() {
	// set up a logrus instance
	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)
	log.SetFormatter(&logrus.TextFormatter{
		DisableSorting:   true,
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		TimestampFormat:  "1504",
	})

	// create the hook
	hookConfig := logrusfluenthook.Config{
		Host:    "127.0.0.1",
		Port:    24224,
		BaseTag: "base_tag",
		Levels: []string{
			"trace",
			"debug",
			"info",
			"warn",
			"error",
			"fatal",
			"panic",
		},
	}
	hook, err := logrusfluenthook.New(&hookConfig)

	// attach the hook/handle any errors
	if err != nil {
		log.Panic(err)
	} else {
		log.AddHook(hook)
	}

	// basic demo
	log.Info("hello from stdout (to fluentd)")

	// demo with `data` fields + custom tag
	log.WithFields(logrus.Fields{
		"this_is": "an example of additional data attachments",
		"tag":     "customise-me",
	}).Info("hello from stdout (to fluentd)")

	// demo with `caller` fields
	log.SetReportCaller(true)
	log.WithFields(logrus.Fields{
		"this_is": "an example of adding caller information",
		"tag":     "with-caller-info",
	}).Debug("hello from stdout (to fluentd)")
}
