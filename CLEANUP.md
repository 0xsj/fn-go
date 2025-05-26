# Commands to fix the port conflict

# 1. First, stop any existing Docker containers

docker compose down
docker ps -a | grep nats | awk '{print $1}' | xargs -r docker stop
docker ps -a | grep nats | awk '{print $1}' | xargs -r docker rm

# 2. Kill any local NATS server processes

sudo lsof -ti:4222 | xargs -r kill -9
sudo lsof -ti:8222 | xargs -r kill -9

# 3. Check if ports are free now

lsof -i :4222
lsof -i :8222

# 4. Start infrastructure

make infra-up
