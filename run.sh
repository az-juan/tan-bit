#!/usr/bin/env bash
set -euo pipefail

ENV_FILE=example.env
DEV_FILE=compose.yaml
PROD_FILE=compose.prod.yaml
export $(grep -v '^#' "$ENV_FILE" | xargs)

if ! command -v gum &>/dev/null; then
  echo "gum no instalado!"
  echo "Visita: https://github.com/charmbracelet/gum"
  exit 1
fi

if ! command -v docker &>/dev/null; then
  gum style --foreground 110 --bold "docker no instalado!, por favor, instalar docker y volver a intentar"
  exit 1
fi

if ! command -v sqlc &>/dev/null; then
  gum style --foreground 110 --bold "sqlc no instalado!, por favor, instalar sqlc y volver a intentar"
  exit 1
fi

if ! command -v atlas &>/dev/null; then
  gum style --foreground 110 --bold "atlas no instalado!, por favor, instalar atlas y volver a intentar"
  exit 1
fi

if ! command -v templ &>/dev/null; then
  gum style --foreground 110 --bold "templ no instalado!, por favor, instalar templ y volver a intentar"
  exit 1
fi

if ! command -v figlet &>/dev/null; then
  gum style --foreground 110 --bold "figlet no instalado!, por favor, instalar figlet y volver a intentar"
  exit 1
fi

if ! command -v go &>/dev/null; then
  gum style --foreground 110 --bold "go no instalado!, por favor, instalar go y volver a intentar"
  exit 1
fi

# MENU
echo "tan-bit" | figlet -f larry3d | gum style --foreground 110

INPUT=$(gum choose "Ejecutar" "Ejecutar tests" "Detener" "Limpiar" "Limpiar cache")

case "$INPUT" in
  "Ejecutar")
    sqlc generate
    templ generate
    docker compose --env-file "$ENV_FILE" -f "$DEV_FILE" -f "$PROD_FILE" up -d 
  ;;

  "Ejecutar tests")
    sqlc generate
    templ generate
    docker compose --env-file "$ENV_FILE" -f "$DEV_FILE" -f "$PROD_FILE" up -d
    atlas migrate apply --dir "file://db/migrations" --url "postgres://"$PG_USER":"$PG_PASSWD"@localhost:5432/"$DB_NAME"?sslmode=disable"
    go test -v ./tests/
  ;;
  "Detener")
    docker compose --env-file "$ENV_FILE" down
  ;;
  "Limpiar")
    docker compose --env-file "$ENV_FILE" down -v --rmi local
  ;;
  "Limpiar cache")
    docker compose --env-file "$ENV_FILE" down -v --rmi local
    docker system prune
  ;;
  *)
    echo "error"
    exit 1
  ;;
esac
