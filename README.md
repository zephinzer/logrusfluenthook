# Logrus Fluent Hook
This library exports a Logrus hook which enables streaming logs to FluentD when added to a Logrus instance.

[![Build Status](https://travis-ci.org/zephinzer/logrusfluenthook.svg?branch=master)](https://travis-ci.org/zephinzer/logrusfluenthook)

# Example Usage

> A full working application that uses this library can be found in [`./cmd/example/main.go`](./cmd/example/main.go).

## Importing

```go
package main

import (
	// ...
  logrusfluenthook "github.com/zephinzer/logrusfluenthook/hook"
  // ...
)
```

## Configuring the Hook

```go
hookConfig := logrusfluenthook.Config{
		Host:    "127.0.0.1",
		Port:    24224,
		BaseTag: "myapplication",
		Levels: []string{
			"trace",
			"debug",
			"info",
			"warn",
			"error",
			"fatal",
			"panic",
		},
		FieldMap: map[string]string{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:   "@msg",
			logrusfluenthook.FieldKeyData:      "@data",
			logrus.FieldKeyLevel:     "@level",
			logrus.FieldKeyFile:      "@file",
			logrus.FieldKeyFunc:      "@func",
		},
		TimestampFormat: time.RFC3339,
	}
```

## Adding the Hook

```go
	logger := logrus.New()
	// the setting of formatter is necessary if you'd like
	// to have the same logs in stdout and in fluentd, if
	// that doesn't matter, you can ignore this next code block
	logger.SetFormatter(&logrus.JSONFormatter{
		DataKey: "@data",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyFile:  "@file",
			logrus.FieldKeyFunc:  "@func",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@msg",
		},
	})
	// ... other customisations of logger ...
	hook, err := logrusfluenthook.New(&hookConfig)
  if err != nil {
		logger.Panic(err)
	} else {
		// note that the fluentd service has to be up and reachable
		// by your application at this point in time, write your
		// own retry/error handler
		logger.AddHook(hook)
	}
```

## Using the Hook

### Basic

```go
// ...
logrus.Info("hello from stdout (to fluentd)")
```

#### Example FluentD Output

```
2019-06-09 14:56:57.683867999 +0000 base_tag.info: {"@timestamp":"2019-06-09T22:56:57+08:00","@msg":"hello from stdout (to fluentd)","@data":{},"@level":"info"}
```

### With Fields

```go
// ...
logrus.WithFields(logrus.Fields{
  "this_is": "an example of additional data attachments",
  "tag":     "customise-me",
}).Info("hello from stdout (to fluentd)")
```

#### Example FluentD Output

```
2019-06-09 14:56:57.685854295 +0000 base_tag.customise-me: {"@timestamp":"2019-06-09T22:56:57+08:00","@msg":"hello from stdout (to fluentd)","@data":{"this_is":"an example of additional data attachments","tag":"customise-me"},"@level":"info"}
```

### With Caller Reporting

```go
// ...
logrus.SetReportCaller(true)
logrus.WithFields(logrus.Fields{
  "this_is": "an example of adding caller information",
  "tag":     "with-caller-info",
}).Debug("hello from stdout (to fluentd)")
```

#### Example FluentD Output

```
2019-06-09 14:56:57.687134165 +0000 base_tag.with-caller-info: {"@msg":"hello from stdout (to fluentd)","@data":{"this_is":"an example of adding caller information","tag":"with-caller-info"},"@level":"debug","@file":"/<REDACTED>/logrusfluenthook/cmd/example/main.go:59","@func":"main.main","@timestamp":"2019-06-09T22:56:57+08:00"}
```

# Development

## Start the FluentD Service
Run `make setup`. This should use Docker Compose to spin up a FluentD service with its ports exposed to your local machine at `127.0.0.1:24224`.

To stop it, use `make teardown`.

## Start Developing
The code is in [`./hook`](./hook) and the example application is at [`./cmd`](./cmd). To check if the example works, you can run `make run_example` from the root directory.

## Testing

Tests can be run using `make run_tests` from the project root directory.

Note that the tests include an integration test which can be found in the [`./test](./test) directory which requires Docker to run (it spins up its own FluentD instance and compares the output/streamed logs).

# Licensing
This repository and the code within is licensed under the MIT license. See [LICENSE](./LICENSE) for the full text.
