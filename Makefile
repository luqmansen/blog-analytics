
GOPATH:=$(shell go env GOPATH)

.PHONY: dev
dev:
	go run main.go


.PHONY: build
	go build -o web-analytics *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t web-analytics:latest