.PHONY: proto docker up down test

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/*/*.proto

docker:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

test:
	go test ./...