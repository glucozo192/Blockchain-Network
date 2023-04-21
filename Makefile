COMPOSE_FILE := ./developments/docker-compose.yml
gen-sql:
	docker compose -f ${COMPOSE_FILE} up generate_sqlc --build
migrate:
	docker compose -f ${COMPOSE_FILE} up migrate
gen-layer:
	go run ./tools/gen-layer/.
	go fmt ./internal