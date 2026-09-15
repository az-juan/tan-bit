APP_NAME := tan-bit
ENV_FILE := example.env

include $(ENV_FILE)
.PHONY: all run generate migrate apply status test stop clean clean_cache

all: run

run: generate
	@docker compose --env-file $(ENV_FILE) up -d

generate:
	@sqlc generate
	@templ generate

migrate:
	@test -n "$(name)" || (echo "Uso: make migrate name=nombre" && exit 1)
	@docker compose --env-file $(ENV_FILE) exec app atlas migrate diff "$(name)" --dir "file://db/migrations" --to \
	"file://db/schema/schema.sql" --dev-url $(DB_URL)

apply:
	@docker compose --env-file $(ENV_FILE) exec app atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"

status:
	@docker compose --env-file $(ENV_FILE) exec app atlas migrate status --dir "file://db/migrations" --url "$(DB_URL)"

test:
	@docker compose --env-file $(ENV_FILE) exec app go test -v .

stop:
	@docker compose --env-file $(ENV_FILE) down

clean:
	@docker compose --env-file $(ENV_FILE) down --volumes --rmi local
	@rm -rf tmp

clean_cache: clean
	@docker system prune
