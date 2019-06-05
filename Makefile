example:
	@$(MAKE) example_mac
	@$(MAKE) example_linux

example_mac:
	@GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
		go build -ldflags '-extldflags "-static"' -o ./bin/example-darwin-amd64 ./cmd/example

example_linux:
	@GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags '-extldflags "-static"' -o ./bin/example-linux-amd64 ./cmd/example

run_example: example
	./bin/example

run_tests:
	@go test -v -cover -coverprofile=c.out ./...

dep:
	@go mod vendor
	@go mod download
deps: # alias for `dep`
	@$(MAKE) dep

setup:
	@docker-compose up -d

teardown:
	@docker-compose down
