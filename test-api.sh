#!/usr/bin/env bash

# Bash Strict Mode
set -euo pipefail

ENV_FILE=example.env
COMPOSE=(docker compose --env-file "$ENV_FILE" -f compose.prod.yaml)
HOST=http://localhost:8080

# Checkear que estan las dependencias instaladas
for dep in docker sqlc go hurl; do
  if ! command -v "$dep" &>/dev/null; then
    echo "$dep no instalado!, por favor, instalar $dep y volver a intentar"
    exit 1
  fi
done

# Exporta las credenciales para que go test (en anfitrion) las lea
set -a
source "$ENV_FILE"
set +a

# MENU
echo
echo "tan-bit"

while true; do
  echo
  echo -e "¿Iniciar servidor y ejecutar tests? [Y/n]\n"n" para bajar Docker\n"
  read -r INPUT
  echo
  INPUT="${INPUT:-Y}"
  if [[ "$INPUT" =~ ^[yY]$ ]]; then
    sqlc generate
    # Bajamos volumenes previos
    "${COMPOSE[@]}" down -v
    # --build: la imagen incluye el codigo/sqlc recien generado
    # --wait: no sigue hasta que ambos contenedores esten healthy
    "${COMPOSE[@]}" up -d --build --wait
    # Sobreescibe variables para que el anfitrion pueda usar ./db/connection/connection.go
    DB_HOST=localhost DB_PORT=5433 go test -count=1 -v ./tests
    hurl --test --variable host="$HOST" tests/hurl_test.hurl
    echo "Servidor de produccion corriendo en $HOST"
  else
    "${COMPOSE[@]}" down -v --rmi local
    docker system prune
    echo "Saliendo..."
    break
  fi
done
