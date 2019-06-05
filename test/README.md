# Logrus Fluent Hook Tests

The tests in here are **not** the unit tests. The unit tests are stored with the source code with the filename postfixed with a `_test` label.

These tests are for manual testing and also to demonstrate how this hook might be used.

## Pre-Requisites

You'll need Docker to run the test environment.

## Get Started

Run the following commands from this directory to run the integration tests:

```sh
# manual run
go test -v ./...

# using Makefile
make
```

## What's Happening

The integration test will setup a FluentD service using Docker, create the logger, add our hook, and make a sample call to the logger. The log received by the FluentD service and the log sent to the output will then be compared for similarity.
