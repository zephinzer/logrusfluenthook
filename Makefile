example:
	@$(MAKE) example_mac
	@$(MAKE) example_linux
	@$(MAKE) example_windows

example_mac: dep
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
		go build -ldflags '-extldflags "-static"' -o ./bin/example-darwin-amd64 ./cmd/example

example_linux: dep
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags '-extldflags "-static"' -o ./bin/example-linux-amd64 ./cmd/example

example_windows: dep
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=386 \
		go build -ldflags '-extldflags "-static"' -o ./bin/example-windows-386.exe ./cmd/example

run_example: example
	./bin/example

test_integration: dep
	@go test -v -cover -coverprofile=c.out ./test

test_unit: dep
	@go test -v -cover -coverprofile=c.out ./hook


dep:
	@go mod vendor
	@go mod download
deps: # alias for `dep`
	@$(MAKE) dep

setup:
	@docker-compose up -d

teardown:
	@docker-compose down
