auth-service:
	@echo	"running auth...."
	@go run ./s



# Makefile for full Docker Compose lifecycle management

DOCKER_COMPOSE_FILE=docker-compose.yaml

.PHONY: down clean prune rebuild all stop remove-images remove-containers remove-volumes

## Stop and remove all containers, networks, and volumes
down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans

## Stop all running containers
stop:
	docker stop $$(docker ps -q) || true

## Remove all stopped containers
remove-containers:
	docker rm $$(docker ps -a -q) || true

## Remove all Docker images
remove-images:
	docker rmi -f $$(docker images -q) || true

## Remove all Docker volumes
remove-volumes:
	docker volume rm $$(docker volume ls -q) || true

## System prune (removes unused containers, networks, images, and volumes)
prune:
	docker system prune -a --volumes -f

## Rebuild containers from scratch
rebuild: down
	docker-compose -f $(DOCKER_COMPOSE_FILE) build
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

## Full reset: stop, remove, clean, rebuild
all: stop remove-containers remove-images remove-volumes rebuild
