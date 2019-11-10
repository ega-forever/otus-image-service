# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=service

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
build_docker:
	docker build -t egorzuev/image_crop_service:1.0.1 .
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
