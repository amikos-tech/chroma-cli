#generate:
#	sh ./gen_api_v3.sh

.PHONY: build
build:
	go build -ldflags "-X 'main.Version=1.0.1-$$(git log -1 --format="%h")' -X 'main.BuildDate=$$(date +%Y-%m-%d)'"

.PHONY: install
install:
	go install -ldflags "-X 'main.Version=1.0.1-$$(git log -1 --format="%h")' -X 'main.BuildDate=$$(date +%Y-%m-%d)'"

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

.PHONY: clean-lint-cache
clean-lint-cache:
	golangci-lint cache clean