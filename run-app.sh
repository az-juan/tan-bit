#!/usr/bin/env bash

# Bash Strict Mode
set -euo pipefail

ENV_FILE=example.env
COMPOSE=(docker compose --env-file "$ENV_FILE" -f compose.prod.yaml)
HOST=http://localhost:8080

# Checkear que estan las dependencias instaladas
for dep in docker sqlc go hurl templ; do
  if ! command -v "$dep" &>/dev/null; then
    echo "$dep no instalado!, por favor, instalar $dep y volver a intentar"
    exit 1
  fi
done

# Exporta las credenciales para que go test (en anfitrion) las lea
# y para que el compose las pueda interpolar
set -a
source "$ENV_FILE"
set +a

# MENU
echo
echo "tan-bit"

while true; do
  echo
  echo -e "¿Iniciar servidor? [Y/n]\n"n" para bajar Docker\n"
  read -r INPUT
  echo
  INPUT="${INPUT:-Y}"
  if [[ "$INPUT" =~ ^[yY]$ ]]; then
    sqlc generate
    # Compila los .templ segun la version declarada en go.mod, no la del anfitrion
    go run "github.com/a-h/templ/cmd/templ@$(go list -m -f '{{.Version}}' github.com/a-h/templ)" generate
    # Bajamos volumenes previos
    "${COMPOSE[@]}" down -v
    # --build: la imagen incluye el codigo/sqlc recien generado
    # --wait: no sigue hasta que ambos contenedores esten healthy
    # directorio para fotos de articulos a publicar
    "${COMPOSE[@]}" up -d --build --wait
    # --file-root para que hurl pueda ver las imagenes de static/assets.
    hurl --test --file-root . --variable host="$HOST" tests/requests.hurl
    echo "Servidor de produccion corriendo en $HOST"
  else
    "${COMPOSE[@]}" down -v --rmi local
    docker system prune
    echo "Saliendo..."
    break
  fi
done
