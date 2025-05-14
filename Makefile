# Variables
SERVICES := user-service auth-service entity-service incident-service location-service notification-service monitoring-service chat-service

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

# Docker development environment
.PHONY: dev-up
dev-up:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

.PHONY: dev-down
dev-down:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

.PHONY: dev-logs
dev-logs:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f

# Rebuild and restart specific service
.PHONY: dev-restart
dev-restart:
	@read -p "Enter service name: " service; \
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build $$service

# Database management
.PHONY: db-shell
db-shell:
	docker exec -it fn-mysql mysql -u$(DB_USER) -p$(DB_PASSWORD)

.PHONY: db-reset
db-reset:
	docker exec -it fn-mysql mysql -uroot -p$(DB_ROOT_PASSWORD) -e "DROP DATABASE IF EXISTS auth_service; DROP DATABASE IF EXISTS user_service; DROP DATABASE IF EXISTS entity_service; DROP DATABASE IF EXISTS incident_service; DROP DATABASE IF EXISTS location_service; DROP DATABASE IF EXISTS monitoring_service; DROP DATABASE IF EXISTS notification_service; DROP DATABASE IF EXISTS chat_service;"
	docker exec -it fn-mysql bash -c "cd /docker-entrypoint-initdb.d && mysql -uroot -p$(DB_ROOT_PASSWORD) < init-databases.sql"

# Deploy to production
.PHONY: prod-build
prod-build:
	docker-compose build

.PHONY: prod-up
prod-up:
	docker-compose up -d

.PHONY: prod-down
prod-down:
	docker-compose down


# Kubernetes commands
.PHONY: k8s-apply
k8s-apply:
	kubectl apply -f infra/k8s/

.PHONY: k8s-delete
k8s-delete:
	kubectl delete -f infra/k8s/

# Observability commands
.PHONY: obs-up
obs-up:
	docker-compose -f infra/observability/docker-compose.yml up -d

.PHONY: obs-down
obs-down:
	docker-compose -f infra/observability/docker-compose.yml down

# Build commands for individual services
.PHONY: build-all
build-all:
	for service in $(SERVICES); do \
		echo "Building $$service..."; \
		go build -o bin/$$service ./services/$$service/cmd/server; \
	done

.PHONY: $(addprefix build-,$(SERVICES))
$(addprefix build-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/build-//'); \
	echo "Building $$service..."; \
	go build -o bin/$$service ./services/$$service/cmd/server

# Run commands for individual services
.PHONY: run-all
run-all:
	nats-server & \
	for service in $(SERVICES); do \
		echo "Starting $$service..."; \
		go run ./services/$$service/cmd/server/main.go & \
	done; \
	cd gateway && go run cmd/server/main.go

.PHONY: $(addprefix run-,$(SERVICES))
$(addprefix run-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/run-//'); \
	echo "Starting $$service..."; \
	cd services/$$service && go run cmd/server/main.go

.PHONY: run-gateway
run-gateway:
	cd gateway && go run main.go

.PHONY: run-local
run-local:
	nats-server & \
	cd services/user-service && go run cmd/server/main.go & \
	cd services/auth-service && go run cmd/server/main.go & \
	cd gateway && go run main.go

# Docker build commands for individual services
.PHONY: $(addprefix docker-build-,$(SERVICES))
$(addprefix docker-build-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/docker-build-//'); \
	echo "Building Docker image for $$service..."; \
	docker build -t $$service:latest -f services/$$service/Dockerfile --build-arg SERVICE_NAME=$$service .

# Testing
.PHONY: test
test:
	go test ./pkg/...
	for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd services/$$service && go test ./... && cd ../..; \
	done
	cd gateway && go test ./... && cd ..

.PHONY: test-coverage
test-coverage:
	go test ./pkg/... -coverprofile=pkg_coverage.out
	for service in $(SERVICES); do \
		echo "Testing $$service with coverage..."; \
		cd services/$$service && go test ./... -coverprofile=../../$$service_coverage.out && cd ../..; \
	done
	cd gateway && go test ./... -coverprofile=../gateway_coverage.out && cd ..
	echo "Merging coverage reports..."
	# You might need to implement custom merge logic here

.PHONY: $(addprefix test-,$(SERVICES))
$(addprefix test-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/test-//'); \
	echo "Testing $$service..."; \
	cd services/$$service && go test ./...

# Linting
.PHONY: lint
lint:
	@echo "Linting pkg directory..."
	cd pkg && golangci-lint run ./...
	@echo "Linting gateway..."
	cd gateway && golangci-lint run ./...
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd services/$$service && golangci-lint run ./... && cd ../..; \
	done

.PHONY: $(addprefix lint-,$(SERVICES))
$(addprefix lint-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/lint-//'); \
	echo "Linting $$service..."; \
	cd services/$$service && golangci-lint run ./...

.PHONY: lint-pkg
lint-pkg:
	@echo "Linting pkg directory..."
	cd pkg && golangci-lint run ./...

.PHONY: lint-gateway
lint-gateway:
	@echo "Linting gateway..."
	cd gateway && golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	@echo "Fixing lint issues in pkg directory..."
	cd pkg && golangci-lint run --fix ./...
	@echo "Fixing lint issues in gateway..."
	cd gateway && golangci-lint run --fix ./...
	@for service in $(SERVICES); do \
		echo "Fixing lint issues in $$service..."; \
		cd services/$$service && golangci-lint run --fix ./... && cd ../..; \
	done

.PHONY: $(addprefix lint-fix-,$(SERVICES))
$(addprefix lint-fix-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/lint-fix-//'); \
	echo "Fixing lint issues in $$service..."; \
	cd services/$$service && golangci-lint run --fix ./...

# Migration commands
.PHONY: $(addprefix migrate-up-,$(SERVICES))
$(addprefix migrate-up-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/migrate-up-//'); \
	echo "Running migrations UP for $$service..."; \
	migrate -path services/$$service/migrations -database "mysql://user:password@tcp(localhost:3306)/$$service" up

.PHONY: $(addprefix migrate-down-,$(SERVICES))
$(addprefix migrate-down-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/migrate-down-//'); \
	echo "Running migrations DOWN for $$service..."; \
	migrate -path services/$$service/migrations -database "mysql://user:password@tcp(localhost:3306)/$$service" down

.PHONY: $(addprefix migrate-create-,$(SERVICES))
$(addprefix migrate-create-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/migrate-create-//'); \
	read -p "Migration name: " name; \
	echo "Creating migration for $$service: $$name"; \
	migrate create -ext sql -dir services/$$service/migrations -seq $$name

# Clean
.PHONY: clean
clean:
	go clean
	docker-compose down -v
	rm -rf bin/
	docker system prune -f

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  Docker:"
	@echo "    docker-build         - Build all Docker images"
	@echo "    docker-up            - Start all containers with docker-compose"
	@echo "    docker-down          - Stop all containers with docker-compose"
	@echo "    docker-build-[service] - Build Docker image for specific service"
	@echo ""
	@echo "  Kubernetes:"
	@echo "    k8s-apply            - Apply all Kubernetes manifests"
	@echo "    k8s-delete           - Delete all Kubernetes manifests"
	@echo ""
	@echo "  Observability:"
	@echo "    obs-up               - Start observability stack (Grafana, Prometheus)"
	@echo "    obs-down             - Stop observability stack"
	@echo ""
	@echo "  Build:"
	@echo "    build-all            - Build all services"
	@echo "    build-[service]      - Build specific service"
	@echo ""
	@echo "  Run:"
	@echo "    run-all              - Run all services"
	@echo "    run-[service]        - Run specific service"
	@echo "    run-gateway          - Run API gateway"
	@echo "    run-local            - Run minimal set of services for local development"
	@echo ""
	@echo "  Testing:"
	@echo "    test                 - Run all tests"
	@echo "    test-coverage        - Run all tests with coverage report"
	@echo "    test-[service]       - Run tests for specific service"
	@echo ""
	@echo "  Linting:"
	@echo "    lint                 - Run linter on all code"
	@echo "    lint-pkg             - Run linter on pkg directory"
	@echo "    lint-gateway         - Run linter on gateway"
	@echo "    lint-[service]       - Run linter on specific service"
	@echo "    lint-fix             - Fix linting issues in all code"
	@echo "    lint-fix-[service]   - Fix linting issues in specific service"
	@echo ""
	@echo "  Migrations:"
	@echo "    migrate-up-[service]   - Run migrations up for specific service"
	@echo "    migrate-down-[service] - Run migrations down for specific service"
	@echo "    migrate-create-[service] - Create new migration for specific service"
	@echo ""
	@echo "  Other:"
	@echo "    clean                - Clean up build artifacts and Docker resources"
	@echo "    help                 - Show this help message"
	@echo ""
	@echo "Available services: $(SERVICES)"