# Makefile for Go KrakenService

IMAGE_NAME=krakenservice
PORT=8080

.PHONY: build docker-build docker-run docker-stop docker-clean

build:
	go build -o krakenservice ./cmd/main.go

docker-build:
	docker build -t $(IMAGE_NAME) .

docker-run:
	docker run --rm -p $(PORT):8080 --name $(IMAGE_NAME) $(IMAGE_NAME)

docker-stop:
	docker stop $(IMAGE_NAME) || true

docker-clean:
	docker rmi $(IMAGE_NAME) || true
	rm -f krakenservice
