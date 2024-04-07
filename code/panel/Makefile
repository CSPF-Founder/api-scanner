DOCKER_STACK = api-scanner

setup: create_dirs cert
	@echo "First time setup done(dont call this agin)!"

create_dirs:
	@echo "First time setup..."
	mkdir -p /app/panel/certs/
	mkdir -p /app/panel/frontend/external/
	mkdir -p /app/data/temp_uploads/
	mkdir -p /app/scanner/local_temp/
	mkdir -p /app/data/work_dir/
	mkdir -p /app/bin/
	mkdir -p /app/manager/local_temp/
	cp -a /vagrant/code/panel/frontend/external/images/. /app/panel/frontend/external/images/

cert:
	@echo "Generating SSL Certificate..."
	openssl req -x509 -sha256 -nodes -days 4096 -newkey rsa:2048 -keyout /app/panel/certs/panel.key -out /app/panel/certs/panel.crt -subj "/C=US/ST=Cysecurity/L=Cysecurity/O=Dis/CN=apiscanner"
	@echo "SSL Generated"

up:
	@echo "Starting containers..."
	docker compose -p ${DOCKER_STACK} up -d
	@echo "Containers started!"

#up_build: starts all containers in the background and recreates them
up_build:
	@echo "Stopping containers..."
	docker compose -p ${DOCKER_STACK} down
	@echo "Building (when required) and Starting containers..."
	docker compose -p ${DOCKER_STACK} up -d --build
	@echo "Docker images build and containers started!"

# down: stops all containers
down:
	@echo "Stopping containers..."
	docker compose -p ${DOCKER_STACK} down
	@echo "Containers stopped!"
