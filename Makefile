example:
	@$(MAKE) example_mac

example_mac:
	@GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
		go build -ldflags '-extldflags "-static"' -o ./bin/example ./cmd/example

run_example: example
	./bin/example