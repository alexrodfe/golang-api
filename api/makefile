.PHONY: build test

test:
	go test `go list ./...` --cover -p=8

# build all docker images
build:
	docker-compose up --build
	