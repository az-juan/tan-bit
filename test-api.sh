#!/usr/bin/env bash
set -euo pipefail

ENV_FILE=example.env
DEV_FILE=compose.yaml
PROD_FILE=compose.prod.yaml

export $(grep -v '^#' "$ENV_FILE" | xargs)

# Checkear que esten las dependencias instaladas
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
  gum style --foreground 111 --bold "sqlc no instalado!, por favor, instalar sqlc y volver a intentar"
  exit 1
fi

if ! command -v atlas &>/dev/null; then
  gum style --foreground 112 --bold "atlas no instalado!, por favor, instalar atlas y volver a intentar"
  exit 1
fi

if ! command -v templ &>/dev/null; then
  gum style --foreground 113 --bold "templ no instalado!, por favor, instalar templ y volver a intentar"
  exit 1
fi

if ! command -v figlet &>/dev/null; then
  gum style --foreground 114 --bold "figlet no instalado!, por favor, instalar figlet y volver a intentar"
  exit 1
fi

if ! command -v go &>/dev/null; then
  gum style --foreground 115 --bold "go no instalado!, por favor, instalar go y volver a intentar"
  exit 1
fi

# MENU
echo "tan-bit" | figlet -f larry3d | gum style --foreground 110

while true; do
  INPUT=$(gum choose "Ejecutar" "Detener" "Salir")
  case "$INPUT" in
    "Ejecutar")
      sqlc generate
      templ generate
      docker compose --env-file "$ENV_FILE" -f "$DEV_FILE" -f "$PROD_FILE" up -d
      atlas migrate apply --dir "file://db/migrations" --url "postgres://"$PG_USER":"$PG_PASSWD"@localhost:5432/"$DB_NAME"?sslmode=disable"
      go test -v ./tests
      chmod u+x ./tests/curl_test.sh
      ./tests/curl_test.sh
    ;;
    "Detener")
      docker compose --env-file "$ENV_FILE" down
    ;;
    "Salir")
      docker compose --env-file "$ENV_FILE" down -v --rmi local
      docker system prune
      echo
      gum style --foreground 115 "saliendo..."
      break
    ;;
    *)
      echo "error"
      exit 1
    ;;
  esac
done
