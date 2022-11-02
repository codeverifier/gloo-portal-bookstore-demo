IMAGE_REPO ?= kasunt
VERSION ?= 1.0.0

.PHONY: docker-local build
all: build docker-local

mod-download:
	go mod download

build: mod-download
	@CGO_ENABLED=0 GOOS=linux go build -a --ldflags '-extldflags "-static"' \
		-installsuffix cgo \
		-o dist/main

docker-build: build
	docker build -t $(IMAGE_REPO)/gloo-portal-bookstore-demo:$(VERSION) .

docker-push:
	docker push $(IMAGE_REPO)/gloo-portal-bookstore-demo:$(VERSION)

clean:
	@rm -fr dist/main