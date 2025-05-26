# FN-GO Setup Guide

## Quick Start

### Option 1: Infrastructure Only (Recommended for Development)

Start only the infrastructure services while developing your Go services locally:

```bash
# Start infrastructure (MySQL, NATS, Prometheus, Grafana, Redis)
make infra-up

# Or equivalently:
make dev-local
```

This will give you:

- **MySQL**: `localhost:3306` (user: `appuser`, password: `apppassword`)
- **NATS**: `localhost:4222` (monitoring: `http://localhost:8222`)
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin/admin)
- **Redis**: `localhost:6379`

Then run your services locally:

```bash
# Run specific services
make run-auth-service
make run-user-service
make run-gateway

# Or run multiple services
make run-local
```

### Option 2: Full Docker Stack

Start everything in Docker:

```bash
# Build and start all services
make docker-up

# Or for development with hot reload
make dev-up
```

## Database Management

```bash
# Check database status
make db-status

# Connect to database
make db-shell          # As appuser
make db-root-shell     # As root

# Reset all databases
make db-reset

# Initialize databases manually
make db-init
```

## Environment Setup

1. Copy the environment template:

```bash
cp .env.example .env
```

2. For infrastructure-only development, you can use:

```bash
cp .env.infrastructure .env
```

## Troubleshooting

### Database Init Issues

If databases aren't being created:

1. Check the init script exists:

```bash
ls -la infra/db/init/init-databases.sql
```

2. Check MySQL logs:

```bash
docker logs fn-mysql
```

3. Manually run the init script:

```bash
make db-init
```

### Service Connection Issues

If services can't connect to infrastructure:

1. Check infrastructure is running:

```bash
docker ps
```

2. Check network connectivity:

```bash
docker exec fn-mysql mysqladmin ping
docker exec fn-nats nats-server --version
```

3. Verify environment variables in your services match the infrastructure ports.

## Development Workflow

1. **Start infrastructure**: `make infra-up`
2. **Run your service locally**: `make run-user-service`
3. **View logs**: `make infra-logs`
4. **Stop when done**: `make infra-down`

## Cleanup

```bash
# Stop and remove containers
make infra-down

# Stop and remove containers + volumes
make infra-clean

# Full cleanup including images
make clean
```
