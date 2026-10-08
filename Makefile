.PHONY: test build
test:
	go test -coverprofile=coverage.out ./...
build:
	./scripts/build.sh
