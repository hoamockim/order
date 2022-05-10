GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test

update:
	go mod tidy

test:
	$(GOTEST) -v ./...

doc:
	cd app/cmd/external/  && swag init

run:
	go run app/cmd/external/main.go

worker_run:
	go run app/cmd/worker/main.go

