PLUGIN_BINARY=systemd-driver
export GO111MODULE=on

default: build

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf ${PLUGIN_BINARY}

build:
	go build -o plugins/${PLUGIN_BINARY} .
