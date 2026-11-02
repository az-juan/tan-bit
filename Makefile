APP_NAME := tan-bit
ENV_FILE := example.env
# Producción por defecto. Para dev: make <target> COMPOSE_FILE=compose.yaml
COMPOSE_FILE ?= compose.prod.yaml
COMPOSE := docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE)
HOST := http://localhost:8080

include $(ENV_FILE)
.PHONY: all run generate api_logs migrate apply status test stop clear clear_cache preCommit

all: run

run: generate
	@$(COMPOSE) up -d --build --wait

generate:
	@sqlc generate

api_logs:
	@$(COMPOSE) logs api

# --- Solo desarrollo: atlas es para dev (COMPOSE_FILE=compose.yaml) ---
migrate:
	@test -n "$(name)" || (echo "Uso: make migrate name=nombre" && exit 1)
	@$(COMPOSE) exec api atlas migrate diff "$(name)" --dir "file://db/migrations" --to \
	"file://db/schema/schema.sql" --dev-url $(DB_URL)

apply:
	@$(COMPOSE) exec api atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"

status:
	@$(COMPOSE) exec api atlas migrate status --dir "file://db/migrations" --url "$(DB_URL)"

# --- Tests en el anfitrión ---
test:
	@set -a && . ./$(ENV_FILE) && set +a && DB_HOST=localhost DB_PORT=5433 go test -count=1 -v ./tests
	@hurl --test --variable host=$(HOST) tests/hurl_test.hurl

stop:
	@$(COMPOSE) down

clear:
	@$(COMPOSE) down --volumes --rmi local

clear_cache: clear
	@docker system prune

preCommit:
	@go mod tidy
	@go fmt ./...
	@go vet ./...
