GOOS = linux
GOARCH = arm64
CGO_ENABLED = 0

.PHONY: all
all: healthcheck-reporter healthcheck-reporter.sha256sum

healthcheck-reporter:
	GOOS="$(GOOS)" \
	GOARCH="$(GOARCH)" \
	CGO_ENABLED="$(CGO_ENABLED)" \
	go build ./cmd/healthcheck-reporter

healthcheck-reporter.sha256sum: healthcheck-reporter
	sha256sum healthcheck-reporter > healthcheck-reporter.sha256sum

.PHONY: clean
clean:
	rm healthcheck-reporter.sha256sum
	rm healthcheck-reporter
