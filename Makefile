# Variables
SERVICES := auth-service user-service entity-service incident-service location-service monitoring-service notification-service chat-service

# Environment setup
.PHONY: setup-env-dev
setup-env-dev:
	@echo "Setting up .env.dev files for all services..."
	@chmod +x setup-env-dev.sh
	@./setup-env-dev.sh

.PHONY: check-env
check-env:
	@echo "Checking environment files:"
	@echo "Root .env file:"
	@if [ -f .env ]; then echo "  ✓ .env exists"; else echo "  ✗ .env missing (copy from .env.example)"; fi
	@echo ""
	@echo "Service .env.dev files:"
	@for service in $(SERVICES); do \
		if [ -f services/$service/.env.dev ]; then \
			echo "  ✓ services/$service/.env.dev"; \
		else \
			echo "  ✗ services/$service/.env.dev missing"; \
		fi; \
	done
	@if [ -f gateway/.env.dev ]; then echo "  ✓ gateway/.env.dev"; else echo "  ✗ gateway/.env.dev missing"; fi

.PHONY: clean-env
clean-env:
	@echo "Removing all .env.dev files..."
	@for service in $(SERVICES); do \
		if [ -f services/$service/.env.dev ]; then \
			rm services/$service/.env.dev; \
			echo "  Removed services/$service/.env.dev"; \
		fi; \
	done
	@if [ -f gateway/.env.dev ]; then rm gateway/.env.dev; echo "  Removed gateway/.env.dev"; fi

# Infrastructure-only commands
.PHONY: infra-up
infra-up:
	@echo "Starting infrastructure services (MySQL, NATS, Prometheus, Grafana, Redis)..."
	docker-compose -f docker-compose.infra.yml up -d

.PHONY: infra-down
infra-down:
	@echo "Stopping infrastructure services..."
	docker-compose -f docker-compose.infra.yml down

.PHONY: infra-logs
infra-logs:
	docker-compose -f docker-compose.infra.yml logs -f

.PHONY: infra-restart
infra-restart:
	@echo "Restarting infrastructure services..."
	docker-compose -f docker-compose.infra.yml restart

.PHONY: infra-clean
infra-clean:
	@echo "Cleaning infrastructure services and volumes..."
	docker-compose -f docker-compose.infra.yml down -v
	docker system prune -f

.PHONY: cleanup-ports
cleanup-ports:
	@echo "Cleaning up port conflicts..."
	@echo "Stopping any existing containers..."
	-docker compose down 2>/dev/null || true
	-docker compose -f docker-compose.infra.yml down 2>/dev/null || true
	@echo "Killing processes on ports 4222 and 8222..."
	-sudo lsof -ti:4222 | xargs -r kill -9 2>/dev/null || true
	-sudo lsof -ti:8222 | xargs -r kill -9 2>/dev/null || true
	@echo "Checking if ports are free..."
	@if lsof -i :4222 >/dev/null 2>&1; then echo "⚠️  Port 4222 still in use"; else echo "✅ Port 4222 is free"; fi
	@if lsof -i :8222 >/dev/null 2>&1; then echo "⚠️  Port 8222 still in use"; else echo "✅ Port 8222 is free"; fi

.PHONY: force-infra-up
force-infra-up: cleanup-ports infra-up

# Database management (works with infra-only setup)
.PHONY: db-shell
db-shell:
	docker exec -it fn-mysql mysql -uappuser -papppassword

.PHONY: db-root-shell
db-root-shell:
	docker exec -it fn-mysql mysql -uroot

.PHONY: db-reset
db-reset:
	@echo "Resetting all databases..."
	docker exec -it fn-mysql mysql -uroot -e "DROP DATABASE IF EXISTS auth_service; DROP DATABASE IF EXISTS user_service; DROP DATABASE IF EXISTS entity_service; DROP DATABASE IF EXISTS incident_service; DROP DATABASE IF EXISTS location_service; DROP DATABASE IF EXISTS monitoring_service; DROP DATABASE IF EXISTS notification_service; DROP DATABASE IF EXISTS chat_service;"
	docker exec -i fn-mysql mysql -uroot < infra/db/init/init-databases.sql

.PHONY: db-init
db-init:
	@echo "Initializing databases..."
	docker exec -i fn-mysql mysql -uroot < infra/db/init/init-databases.sql

.PHONY: db-status
db-status:
	@echo "Checking database status..."
	docker exec -it fn-mysql mysql -uroot -e "SHOW DATABASES;"

.PHONY: db-create-appuser
db-create-appuser:
	@echo "Creating appuser and databases..."
	docker exec -it fn-mysql mysql -uroot -e "CREATE USER IF NOT EXISTS 'appuser'@'%' IDENTIFIED BY 'apppassword';"
	docker exec -it fn-mysql mysql -uroot -e "GRANT ALL PRIVILEGES ON *.* TO 'appuser'@'%' WITH GRANT OPTION;"
	docker exec -it fn-mysql mysql -uroot -e "FLUSH PRIVILEGES;"
	@echo "Creating databases..."
	docker exec -it fn-mysql mysql -uroot -e "CREATE DATABASE IF NOT EXISTS auth_service;"
	docker exec -it fn-mysql mysql -uroot -e "CREATE DATABASE IF NOT EXISTS user_service;"
	docker exec -it fn-mysql mysql -uroot -e "CREATE DATABASE IF NOT EXISTS entity_service;"
	docker exec -it fn-mysql mysql -uroot -e "CREATE DATABASE IF NOT EXISTS incident_service;"
	docker exec -it fn-mysql mysql -uroot -e "CREATE DATABASE IF NOT EXISTS location_service;"
	docker exec -it fn-mysql mysql -uroot -e "CREATE DATABASE IF NOT EXISTS monitoring_service;"
	docker exec -it fn-mysql mysql -uroot -e "CREATE DATABASE IF NOT EXISTS notification_service;"
	docker exec -it fn-mysql mysql -uroot -e "CREATE DATABASE IF NOT EXISTS chat_service;"
	@echo "Granting privileges to appuser..."
	docker exec -it fn-mysql mysql -uroot -e "GRANT ALL PRIVILEGES ON auth_service.* TO 'appuser'@'%';"
	docker exec -it fn-mysql mysql -uroot -e "GRANT ALL PRIVILEGES ON user_service.* TO 'appuser'@'%';"
	docker exec -it fn-mysql mysql -uroot -e "GRANT ALL PRIVILEGES ON entity_service.* TO 'appuser'@'%';"
	docker exec -it fn-mysql mysql -uroot -e "GRANT ALL PRIVILEGES ON incident_service.* TO 'appuser'@'%';"
	docker exec -it fn-mysql mysql -uroot -e "GRANT ALL PRIVILEGES ON location_service.* TO 'appuser'@'%';"
	docker exec -it fn-mysql mysql -uroot -e "GRANT ALL PRIVILEGES ON monitoring_service.* TO 'appuser'@'%';"
	docker exec -it fn-mysql mysql -uroot -e "GRANT ALL PRIVILEGES ON notification_service.* TO 'appuser'@'%';"
	docker exec -it fn-mysql mysql -uroot -e "GRANT ALL PRIVILEGES ON chat_service.* TO 'appuser'@'%';"
	docker exec -it fn-mysql mysql -uroot -e "FLUSH PRIVILEGES;"
	@echo "✅ Setup complete! Testing appuser connection..."
	docker exec -it fn-mysql mysql -uappuser -papppassword -e "SHOW DATABASES;"

.PHONY: db-check-connection
db-check-connection:
	@echo "Checking MySQL container status..."
	@if docker ps | grep -q fn-mysql; then \
		echo "✅ MySQL container is running"; \
		echo "Checking MySQL connectivity..."; \
		docker exec fn-mysql mysqladmin ping 2>/dev/null && echo "✅ MySQL is responding" || echo "❌ MySQL is not responding"; \
	else \
		echo "❌ MySQL container is not running"; \
		echo "Run 'make infra-up' to start infrastructure"; \
	fi

.PHONY: db-debug
db-debug:
	@echo "🔍 MySQL Debug Information:"
	@echo "Container status:"
	@docker ps | grep mysql || echo "No MySQL containers running"
	@echo ""
	@echo "Container logs (last 20 lines):"
	@docker logs --tail 20 fn-mysql 2>/dev/null || echo "Cannot get logs"
	@echo ""
	@echo "Environment variables:"
	@docker exec fn-mysql env | grep MYSQL 2>/dev/null || echo "Cannot get environment"
	@echo ""
	@echo "Trying different connection methods:"
	@echo "1. Root with no password:"
	@docker exec fn-mysql mysql -uroot -e "SELECT 1;" 2>/dev/null && echo "✅ Root with no password works" || echo "❌ Root with no password failed"
	@echo "2. Root with 'rootpassword':"
	@docker exec fn-mysql mysql -uroot -prootpassword -e "SELECT 1;" 2>/dev/null && echo "✅ Root with 'rootpassword' works" || echo "❌ Root with 'rootpassword' failed"
	@echo "3. Appuser with 'apppassword':"
	@docker exec fn-mysql mysql -uappuser -papppassword -e "SELECT 1;" 2>/dev/null && echo "✅ Appuser works" || echo "❌ Appuser failed"

.PHONY: db-reset-container
db-reset-container:
	@echo "🚨 RESETTING MySQL container (this will delete all data!)"
	@read -p "Are you sure? Type 'yes' to continue: " confirm; \
	if [ "$confirm" = "yes" ]; then \
		echo "Stopping and removing MySQL container..."; \
		docker-compose -f docker-compose.infra.yml stop mysql; \
		docker-compose -f docker-compose.infra.yml rm -f mysql; \
		echo "Removing MySQL volume..."; \
		docker volume rm fn-go_mysql-data 2>/dev/null || true; \
		echo "Starting fresh MySQL container..."; \
		docker-compose -f docker-compose.infra.yml up -d mysql; \
		echo "Waiting for MySQL to initialize..."; \
		sleep 10; \
		make db-debug; \
	else \
		echo "Reset cancelled."; \
	fi

# Local development with infrastructure
.PHONY: dev-local
dev-local: infra-up
	@echo "Infrastructure started. Setting up database..."
	@sleep 5
	@make db-create-appuser
	@echo ""
	@echo "🚀 Development environment ready!"
	@echo ""
	@echo "📊 Services available:"
	@echo "  NATS: localhost:4222"
	@echo "  NATS Monitor: http://localhost:8222"
	@echo "  MySQL: localhost:3306"
	@echo "  Prometheus: http://localhost:9090"
	@echo "  Grafana: http://localhost:3000 (admin/admin)"
	@echo "  Redis: localhost:6379"
	@echo ""
	@echo "🔐 Database credentials:"
	@echo "  Root: root / (no password)"
	@echo "  App:  appuser / apppassword"
	@echo ""
	@echo "▶️  Start your services:"
	@echo "  make run-auth-service"
	@echo "  make run-user-service"
	@echo "  make run-gateway"

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

# Docker development environment (full stack)
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

# Run commands for individual services (assumes infrastructure is running)
.PHONY: run-all
run-all:
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
	cd gateway && go run cmd/server/main.go

.PHONY: run-local
run-local:
	cd services/user-service && go run cmd/server/main.go & \
	cd services/auth-service && go run cmd/server/main.go & \
	cd gateway && go run cmd/server/main.go

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

# Migration commands (FIXED VERSION)

.PHONY: migrate-all-up
migrate-all-up:
	@echo "Running migrations UP for all services..."
	@for service in $(SERVICES); do \
		db_name=$$(echo $$service | tr '-' '_'); \
		echo "📊 Running migrations for $$service (database: $$db_name)..."; \
		migrate -path services/$$service/migrations -database "mysql://appuser:apppassword@tcp(localhost:3306)/$$db_name" up || echo "❌ Failed to migrate $$service"; \
		echo ""; \
	done
	@echo "✅ All migrations completed!"

.PHONY: migrate-all-down
migrate-all-down:
	@echo "🚨 Running migrations DOWN for all services (this will rollback changes)..."
	@read -p "Are you sure? Type 'yes' to continue: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		for service in $(SERVICES); do \
			db_name=$$(echo $$service | tr '-' '_'); \
			echo "📊 Rolling back migrations for $$service (database: $$db_name)..."; \
			migrate -path services/$$service/migrations -database "mysql://appuser:apppassword@tcp(localhost:3306)/$$db_name" down || echo "❌ Failed to rollback $$service"; \
			echo ""; \
		done; \
		echo "✅ All rollbacks completed!"; \
	else \
		echo "Rollback cancelled."; \
	fi

.PHONY: migrate-all-status
migrate-all-status: migrate-status


# INDIVIDUAL
.PHONY: $(addprefix migrate-up-,$(SERVICES))
$(addprefix migrate-up-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/migrate-up-//'); \
	db_name=$$(echo $$service | tr '-' '_'); \
	echo "Running migrations UP for $$service..."; \
	migrate -path services/$$service/migrations -database "mysql://appuser:apppassword@tcp(localhost:3306)/$$db_name" up

.PHONY: $(addprefix migrate-down-,$(SERVICES))
$(addprefix migrate-down-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/migrate-down-//'); \
	db_name=$$(echo $$service | tr '-' '_'); \
	echo "Running migrations DOWN for $$service..."; \
	migrate -path services/$$service/migrations -database "mysql://appuser:apppassword@tcp(localhost:3306)/$$db_name" down

.PHONY: $(addprefix migrate-create-,$(SERVICES))
$(addprefix migrate-create-,$(SERVICES)):
	@service=$$(echo $@ | sed 's/migrate-create-//'); \
	read -p "Migration name: " name; \
	echo "Creating migration for $$service: $$name"; \
	migrate create -ext sql -dir services/$$service/migrations -seq $$name

.PHONY: migrate-status
migrate-status:
	@echo "Checking migration status for all services..."
	@for service in $(SERVICES); do \
		db_name=$$(echo $$service | tr '-' '_'); \
		echo "📊 $$service (database: $$db_name):"; \
		migrate -path services/$$service/migrations -database "mysql://appuser:apppassword@tcp(localhost:3306)/$$db_name" version 2>/dev/null || echo "  No migrations applied yet"; \
		echo ""; \
	done

.PHONY: migrate-force
migrate-force:
	@read -p "Enter service name: " service; \
	read -p "Enter version to force: " version; \
	db_name=$$(echo $$service | tr '-' '_'); \
	echo "Forcing migration version $$version for $$service..."; \
	migrate -path services/$$service/migrations -database "mysql://appuser:apppassword@tcp(localhost:3306)/$$db_name" force $$version

# Clean
.PHONY: clean
clean:
	go clean
	docker-compose down -v
	docker-compose -f docker-compose.infrastructure.yml down -v
	rm -rf bin/
	docker system prune -f

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo ""
	@echo "  Environment:"
	@echo "    setup-env-dev        - Create .env.dev files for all services"
	@echo "    check-env            - Check status of environment files"
	@echo "    clean-env            - Remove all .env.dev files"
	@echo ""
	@echo "  Infrastructure (for local development):"
	@echo "    infra-up             - Start infrastructure services only (MySQL, NATS, etc.)"
	@echo "    infra-down           - Stop infrastructure services"
	@echo "    infra-logs           - Show infrastructure logs"
	@echo "    infra-restart        - Restart infrastructure services"
	@echo "    infra-clean          - Stop infrastructure and clean volumes"
	@echo "    cleanup-ports        - Clean up port conflicts (kills processes on 4222, 8222)"
	@echo "    force-infra-up       - Clean ports and start infrastructure"
	@echo "    dev-local            - Start infrastructure and show connection info"
	@echo ""
	@echo "  Database:"
	@echo "    db-check-connection  - Check MySQL container status"
	@echo "    db-debug             - Debug MySQL connection issues"
	@echo "    db-create-appuser    - Create appuser and databases (run this first!)"
	@echo "    db-shell             - Connect to MySQL as appuser"
	@echo "    db-root-shell        - Connect to MySQL as root (no password)"
	@echo "    db-reset             - Reset all databases"
	@echo "    db-reset-container   - Reset entire MySQL container (DESTRUCTIVE)"
	@echo "    db-init              - Initialize databases"
	@echo "    db-status            - Show database status"
	@echo ""
	@echo "  Docker (full stack):"
	@echo "    docker-build         - Build all Docker images"
	@echo "    docker-up            - Start all containers with docker-compose"
	@echo "    docker-down          - Stop all containers with docker-compose"
	@echo "    docker-build-[service] - Build Docker image for specific service"
	@echo ""
	@echo "  Development (full stack):"
	@echo "    dev-up               - Start all services in development mode"
	@echo "    dev-down             - Stop all development services"
	@echo "    dev-logs             - Show development logs"
	@echo "    dev-restart          - Restart specific development service"
	@echo ""
	@echo "  Kubernetes:"
	@echo "    k8s-apply            - Apply all Kubernetes manifests"
	@echo "    k8s-delete           - Delete all Kubernetes manifests"
	@echo ""
	@echo "  Build:"
	@echo "    build-all            - Build all services"
	@echo "    build-[service]      - Build specific service"
	@echo ""
	@echo "  Run (requires infrastructure):"
	@echo "    run-all              - Run all services locally"
	@echo "    run-[service]        - Run specific service locally"
	@echo "    run-gateway          - Run API gateway locally"
	@echo "    run-local            - Run minimal set of services locally"
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
	@echo "    migrate-up-[service]     - Run migrations up for specific service"
	@echo "    migrate-down-[service]   - Run migrations down for specific service"
	@echo "    migrate-create-[service] - Create new migration for specific service"
	@echo "    migrate-status           - Check migration status for all services"
	@echo "    migrate-force            - Force migration version (use with caution)"
	@echo ""
	@echo "  Other:"
	@echo "    clean                - Clean up build artifacts and Docker resources"
	@echo "    help                 - Show this help message"
	@echo ""
	@echo "Available services: $(SERVICES)"