APP_NAME := my-app
DB_URL := postgres://postgres:ayudanoentiendo@localhost:5433/tan_bit?sslmode=disable

# Declaracion de targets
.PHONY: all run generate migrate apply status build test clean

all: build

run:
	@air

# Ejecuta las tareas de generación configuradas por el proyecto
generate:
	@sqlc generate

# @templ generate

# Genera una migración: make migrate name=add_created_at
# No queremos incluirla dentro del build
migrate:
	@test -n "$(name)" || (echo "Uso: make migrate name=nombre" && exit 1)
	@atlas migrate diff "$(name)" --dir "file://db/migrations" --to \
	"file://db/schema/schema.sql" --dev-url "docker://postgres/15/dev?search_path=public"

# Aplica las migraciones pendientes
apply:
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"

# Muestra el estado de las migraciones
status:
	atlas migrate status --dir "file://db/migrations" --url "$(DB_URL)"

# Construye el binario de la aplicación
build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

test:
	@go test ./...

# Limpia los artefactos de construcción
clean:
	@rm -rf tmp

preCommit:
	@go mod tidy
	@go fmt
	@go vet
	@echo "Testing result:"
	@go test