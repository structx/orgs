
build:
	docker build -t structx/orgs:v0.1.0 .

deps:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run