include example.env

APP_NAME := tan-bit
COMPOSE_PROD := compose.prod.yaml
COMPOSE_DEV := compose.yaml

.PHONY: all run generate migrate apply status test clean

all: run

run: generate
	@docker compose --env-file example.env -f $(COMPOSE_DEV) up -d

# Ejecuta las tareas de generación configuradas por el proyecto
generate:
	@sqlc generate
	@templ generate

# Genera una migración: make migrate name=add_created_at
migrate:
	@test -n "$(name)" || (echo "Uso: make migrate name=nombre" && exit 1)
	@docker compose --env-file example.env -f $(COMPOSE_DEV) exec app atlas migrate diff "$(name)" --dir "file://db/migrations" --to \
	"file://db/schema/schema.sql" --dev-url "docker://postgres/15/dev?search_path=public"

# Aplica las migraciones pendientes
apply:
	@docker compose --env-file example.env -f $(COMPOSE_DEV) exec app atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"

# Muestra el estado de las migraciones
status:
	@docker compose --env-file example.env -f $(COMPOSE_DEV) exec app atlas migrate status --dir "file://db/migrations" --url "$(DB_URL)"

test:
	@docker compose --env-file example.env exec app go test -v .

# Limpia los artefactos de construcción
clean:
	@docker compose -f $(COMPOSE_DEV) rm -s -f -v
	@rm -rf tmp
