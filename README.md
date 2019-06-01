# Logrus Fluent Hook
This library exports a Logrus hook which enables streaming logs to FLuentD when added to a Logrus instance.

# Example Usage

> A full working application that uses this library can be found in [`./cmd/example/main.go`](./cmd/example/main.go).

## Importing

```go
package main

import (
  logrusfluenthook "github.com/zephinzer/logrusfluenthook/hook"
  // ...
)
```

## Configuring the Hook

```go
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
```

## Adding the Hook

```go
	hook, err := logrusfluenthook.New(&hookConfig)
  if err != nil {
		logrus.Panic(err)
	} else {
		logrus.AddHook(hook)
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
2019-06-01 18:23:20.035382000 +0000 base_tag.log: {"@timestamp":"2019-06-02T02:23:20+08:00","@msg":"hello from stdout (to fluentd)","@data":{}}
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
2019-05-31 13:04:11.722584000 +0000 base_tag.customise-me: {"timestamp":"2019-05-31T15:04:11+02:00","msg":"hello from stdout (to fluentd)","data":{"tag":"customise-me","this_is":"an example of additional data attachments"}}
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
2019-05-31 13:04:11.739373000 +0000 base_tag.with-caller-info: {"timestamp":"2019-05-31T15:04:11+02:00","msg":"hello from stdout (to fluentd)","data":{"tag":"with-caller-info","this_is":"an example of adding caller information"},"caller":{"file":"/Users/zephinzer/Projects/logrus_fluent_hook/cmd/example/main.go","line":59,"function":"main.main"}}
```

# Development

## Start the FluentD Service
Navigate to `./test` and start the FluentD service using `make start`. This should use Docker Compose to spin up a FluentD service with its ports exposed to your local machine at `127.0.0.1:24224`.

## Start Developing
The code is in [`./hook`](./hook) and the example application is at [`./cmd`](./cmd). To check if the example works, you can run `make run_example` from the root directory.

## Testing
WIP

# Licensing
This repository and the code within is licensed under the MIT license. See [LICENSE](./LICENSE) for the full text.
