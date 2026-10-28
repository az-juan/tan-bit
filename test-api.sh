#!/usr/bin/env bash

# Bash Strict Mode
set -euo pipefail

ENV_FILE=example.env
DEV_FILE=compose.yaml
PROD_FILE=compose.prod.yaml

export $(grep -v '^#' "$ENV_FILE" | xargs)

# Checkear que esten las dependencias instaladas
if ! command -v docker &>/dev/null; then
  echo "Docker no instalado!, por favor, instalar docker y volver a intentar"
  exit 1
fi

if ! command -v sqlc &>/dev/null; then
  echo "sqlc no instalado!, por favor, instalar sqlc y volver a intentar"
  exit 1
fi

if ! command -v templ &>/dev/null; then
  echo "templ no instalado!, por favor, instalar templ y volver a intentar"
  exit 1
fi

if ! command -v go &>/dev/null; then
  echo "Go no instalado!, por favor, instalar go y volver a intentar"
  exit 1
fi

# MENU
echo "tan-bit"

while true; do
  read -r -n 1 -p "¿Iniciar servidor? [Y/n] " INPUT
    echo ""
    INPUT="${INPUT:-Y}"
    if [[ "$response" =~ ^[yY]$ ]]; then
      sqlc generate
      templ generate
      # wait no permite que el script siga sin que los dos contenedores esten arriba
      docker compose --env-file "$ENV_FILE" -f "$DEV_FILE" -f "$PROD_FILE" up -d --wait
      go test -v ./tests
      chmod u+x ./tests/curl_test.sh
      ./tests/curl_test.sh
    else
      docker compose --env-file "$ENV_FILE" down -v --rmi local
      docker system prune
      echo "Saliendo..."
      break
    fi
done
