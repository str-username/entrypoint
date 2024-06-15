### Variables ###
CURRENT_DIRECTORY=$(shell pwd)
COMPOSE_FILE="deploy/local/docker-compose.yaml"
DOCKER_FILE="build/Dockerfile"
DOCKER_CTX="."

# Local: start mongo
start_mongo:
	docker compose -f ${COMPOSE_FILE} up mongodb -d

build_container:
	docker build -f ${DOCKER_FILE} ${DOCKER_CTX}

# Local: cleanup
cleanup:
	docker compose -f ${COMPOSE_FILE} down