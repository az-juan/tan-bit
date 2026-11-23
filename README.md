# Mi primera aplicacion web, TAN-BIT

![Alt text](./images/tanbit_logo.png)

Mercado de hardware/tecnología destinado a darle oportunidad a articulos que aún tienen utilidad entre vecinos.
Hardware de segunda mano y hardware custom/personalizado hecho de forma 'artesanal'.

*DESCENTRALIZADO - CONFIANZA - OPORTUNIDAD*

## Dependencias
+ docker
+ docker compose
+ go 1.26 (versión exacta en go.mod)
+ sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
+ templ: `go install github.com/a-h/templ/cmd/templ@latest`
+ hurl: https://hurl.dev

## Estructura del proyecto
Vista general completa en [`docs/arquitectura.md`](./docs/arquitectura.md).

```
tan-bit.com/tan-bit
├── main.go              # pool DB, rutas, static fs
├── go.mod / go.sum      # go 1.26
├── Dockerfile*          # dev (air) y prod (binario alpine)
├── nginx.conf           # proxy /api/ y /uploads/ (producción)
├── compose.yaml         # desarrollo (api + db)
├── compose.prod.yaml    # producción (api + db + nginx)
├── example.env          # variables DB_*
├── Makefile             # run/test/migrate/apply/status/preCommit/clear
├── run-app.sh           # arranque con sqlc+templ, compose y tests hurl
├── db/
│   ├── schema/          # schema.sql (DDL real de la BD)
│   ├── queries/         # consultas sqlc (-- name: ...)
│   ├── sqlc/            # codegen sqlc
│   ├── migrations/      # migraciones Atlas (dev)
│   └── connection/      # pool pgx
├── logic/               # regla de negocio (parseo, validación, auth, notificaciones, tx)
├── handlers/            # net/http handlers
├── views/               # componentes SSR con templ + htmx
├── static/              # css, assets, htmx.min.js
└── tests/               # go test + escenarios Hurl
```

## Ejecutar la aplicación
```bash
chmod u+x run-app.sh
./run-app.sh` # — bucle: genera, baja volúmenes, `up -d --build --wait`, ejecuta hurl; `n` para salir.
```
Accedader a la página: http://localhost:8080

## Credenciales de prueba (vendedor pre-cargado)
Al iniciar la aplicación con ./run-app.sh, los tests de Hurl crean automáticamente un vendedor con artículos publicados en el catálogo:
+ **Email:** leo@messi.com
+ **contraseña:** campeon2022

## Asegurarse de
+ Estar en el grupo 'docker', en caso de no estarlo, ejecutar: `sudo usermod -aG docker $USER`.
Al ejecutar `groups` tiene que listarse "docker", de lo contrario cierra la sesión actual e inicia otra.
+ Tener puertos 80, 5433 y 8080 libres
+ Iniciar el Docker Daemon con: `sudo systemctl start docker`
+ No tener bloqueador de cookies habilitado
