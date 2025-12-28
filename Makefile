APP_NAME=ecomtechgo
DOCKER_IMAGE=$(APP_NAME)

.PHONY: run
run:
	go run cmd/main.go

.PHONY: build
build:
	go build -o $(APP_NAME) cmd/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: clean
clean:
	rm -f $(APP_NAME)

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run:
	docker run -p 8081:8081 $(DOCKER_IMAGE)
