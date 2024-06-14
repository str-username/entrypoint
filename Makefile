### Variables ###
CURRENT_DIRECTORY=$(shell pwd)
COMPOSE_FILE="deploy/local/docker-compose.yaml"

# Local: start mongo
start_mongo:
	docker compose -f ${COMPOSE_FILE} up mongodb -d

# Local: cleanup
cleanup:
	docker compose -f ${COMPOSE_FILE} down
