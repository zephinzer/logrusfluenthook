# Logrus Fluent Hook Tests

The tests in here are **not** the unit tests. The unit tests are stored with the source code with the filename postfixed with a `_test` label.

These tests are for manual testing and also to demonstrate how this hook might be used.

## Pre-Requisites

You'll need Docker and Docker Compose to run the test environment.

## Get Started

Run the following command to spin up FluentD:

```sh
docker-compose up
```

To background it:

```sh
docker-compose up -d
```

To view the logs when it's backgrounded:

```sh
docker-compose logs -f
```
