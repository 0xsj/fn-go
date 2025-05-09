# Docker build commands
.PHONY: docker-build
docker-build:
	docker-compose build

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

# Kubernetes commands
.PHONY: k8s-apply
k8s-apply:
	kubectl apply -f infra/kubernetes/

.PHONY: k8s-delete
k8s-delete:
	kubectl delete -f infra/kubernetes/

# Local development
.PHONY: run-local
run-local:
	nats-server & \
	cd services/user-service && go run main.go & \
	cd gateway && go run main.go

# Testing
.PHONY: test
test:
	go test ./...

# Clean
.PHONY: clean
clean:
	go clean
	docker-compose down -v