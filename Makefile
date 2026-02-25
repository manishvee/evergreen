CMD_DIR  := ./cmd/evergreen

.PHONY: build
build:
	@mkdir -p ./bin
	go build -o ./bin/evergreen $(CMD_DIR)

.PHONY: dev
dev:
	docker run --rm -it \
		-v .:/workspace/ \
		-w /workspace/ \
		-p 5225:5225 \
		golang:trixie

.PHONY: run
run:
	go run $(CMD_DIR)

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...
