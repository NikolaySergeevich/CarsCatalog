PKG_LIST := $(shell go list ./... | grep -v /vendor/)
PATH := $(PATH):$(GOPATH)/bin

.PHONY: build
build:
	go build -o bin/auto-srv cmd/auto-srv/main.go

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: genid
genid:
	go run cmd/genid/main.go

.PHONY: generate
generate:

	go generate ./...


.PHONY: test
test:
	go test -short ./...

.PHONY: integration
integration:
	go test -race ./...

.PHONY: docker
docker-run:
	docker run --name autoCatalog -e POSTGRES_DB=auto -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d -p "5434:5432" postgres

.PHONY: migrate-up
migrate-up:
	 migrate -source "file://./migrations" -database "postgres://localhost:5434/auto?sslmode=disable&user=postgres&password=postgres" up

.PHONY: migrate-down
migrate-down:
	 migrate -source "file://./migrations" -database "postgres://localhost:5434/auto?sslmode=disable&user=postgres&password=postgres" down
