# Logrus Fluent Hook
This library exports a Logrus hook which enables streaming logs to FLuentD when added to a Logrus instance.

# Example Usage

> A full working application that uses this library can be found in [`./cmd/example/main.go`](./cmd/example/main.go).

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
logrus.Info("hello from stdout (to fluentd)")
```

#### Example FluentD Output

```
2019-05-31 13:04:11.717347000 +0000 base_tag.log: {"data":{"tag":"log"},"timestamp":"2019-05-31T15:04:11+02:00","msg":"hello from stdout (to fluentd)"}
```

### With Fields

```go
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

# Licensing
This repository and the code within is licensed under the MIT license. See [LICENSE](./LICENSE) for the full text.
