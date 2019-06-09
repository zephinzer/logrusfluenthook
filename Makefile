GOOS :=
GOARCH :=
BINEXT :=

ifeq ($(OS),Windows_NT)
	# this is windows
	GOOS = windows

	# get GOARCH
	ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
		GOARCH = amd64
	else
		ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
			GOARCH = amd64
		endif
		ifeq ($(PROCESSOR_ARCHITECTURE),x86)
				GOARCH = 386
		endif
	endif
else
	# get GOOS
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS = darwin
	endif
	# get GOARCH
	UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		GOARCH = amd64
	endif
	ifneq ($(filter %86,$(UNAME_P)),)
		GOARCH = 386
	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		GOARCH = arm
	endif
endif

ifeq ($(GOOS),windows)
	BINEXT = .exe
endif

# use these for dependency installation
dep:
	@go mod vendor
	@go mod download
deps: # alias for `dep`
	@$(MAKE) dep

# use this to setup fluentd locally
setup:
	@cd deployments && docker-compose up -d -V

# use this to shutdown everything
teardown:
	@cd deployments && docker-compose down

# use this to run the example
run_example: example setup
	@printf -- "\n\nwaiting 5 seconds for fluentd to initialise"
	@sleep 1
	@printf -- '.'
	@sleep 1
	@printf -- '.'
	@sleep 1
	@printf -- '.'
	@sleep 1
	@printf -- '.'
	@sleep 1
	@printf -- '.\n'
	@sleep 1
	./bin/example-$(GOOS)-$(GOARCH)$(BINEXT)
	@printf -- "\n\nlogs received by fluentd:\n"
	@docker logs logrus_fluent_hook_fluent_instance | tail -n 3
	@printf -- 'psa: run make teardown once you no longer need fluentd!\n\n'

# use this to run the integration test
test_integration: dep
	@go test -v -cover -coverprofile=c.out ./test

# use this to run the unit tests
test_unit: dep
	@go test -v -cover -coverprofile=c.out ./hook

##################################
# utility functions & more below #
##################################

# creates the example depending on your operating system
example:
	GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags '-extldflags "-static"' -o ./bin/example-$(GOOS)-$(GOARCH)$(BINEXT) ./cmd/example

# creates the example for most macoses
example_mac: dep
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
		go build -ldflags '-extldflags "-static"' -o ./bin/example-darwin-amd64 ./cmd/example

# creates the example for most linuxes
example_linux: dep
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags '-extldflags "-static"' -o ./bin/example-linux-amd64 ./cmd/example

# creates the example for mosts windows
example_windows: dep
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=386 \
		go build -ldflags '-extldflags "-static"' -o ./bin/example-windows-386.exe ./cmd/example
