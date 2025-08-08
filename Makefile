# Makefile for the ArgoPipeline Operator

IMG ?= r-operator:latest

all: build

build:
	go build -o bin/manager main.go

run: generate fmt vet
	go run ./main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

test: generate fmt vet
	go test ./... -coverprofile cover.out

docker-build:
	docker build -t ${IMG} .

docker-push:
	docker push ${IMG}

generate:
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

manifests:
	controller-gen rbac:roleName=manager-role crd paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: all build run fmt vet test docker-build docker-push generate manifests
