# Diagramas de arquitectura — TAN-BIT

## 1. Estructura de módulos (module view)

```mermaid
flowchart TB
  subgraph app["tan-bit.com/tan-bit (Go module)"]
    main["main.go\n(entrypoint, routing, static fs)"]
    handlers["handlers/\nArticuloHandler (ServeHTTP)\nMarketplace/Carrito/Evento/..."]
    views["views/\n*.templ + *_templ.go\nHeader(), Footer()"]
    dbpkg["db/"]
    static["static/\ncss/*.css, assets/*"]
  end
  subgraph dbpkg2["db/"]
    conn["connection/\nOpen()"]
    sqlc["sqlc/\nmodels, queries (generated)"]
    queries["queries/\n*.sql (sqlc input)"]
    schema["schema/\nschema.sql"]
    migrations["migrations/\nAtlas migrations"]
  end
  main --> handlers
  main --> views
  main --> conn
  main --> static
  handlers --> views
  handlers --> sqlc
  conn --> sqlc
  queries -.sqlc generate.-> sqlc
  schema -.Atlas.-> migrations
```

## 2. Componentes y conectores (C&C view, runtime)

```mermaid
flowchart LR
  browser["Browser\n(HTML/CSS)"]
  http["Go HTTP Server\n:8080\nhttp.ServeMux"]
  handler["ArticuloHandler\nhandlers"]
  renderer["Templ renderer\nviews/*_templ.go"]
  sqlc["sqlc queries\npgx pool"]
  pg["PostgreSQL\n:5432"]
  fs["Static file server\n/static/*"]

  browser -- "HTTP GET/POST" --> http
  http --> handler
  http --> fs
  handler --> renderer
  handler --> sqlc
  sqlc -- "SQL (pgx)" --> pg
  renderer -- "HTML" --> browser
  fs -- "CSS/assets" --> browser
```

## 3. Despliegue (allocation view)

```mermaid
flowchart TB
  subgraph compose["Docker Compose"]
    subgraph appc["Container: app (Go)"]
      air["air (hot reload, dev)"]
      goapp["Go binary (prod)"]
      templ["templ-generated code"]
    end
    subgraph dbc["Container: db (PostgreSQL)"]
      pg["PostgreSQL 15"]
      data["Volume: pgdata"]
    end
  end
  host["Host\n:8080 (app)\n:5433 (db dev)"]
  host --> appc
  appc -- "TCP 5432" --> dbc
  dbc --> data
```

## 4. Modelo de datos (ERD)

```mermaid
erDiagram
  USUARIO {
    int id_usuario PK
    varchar apellido
    varchar nombre
    varchar email UK
    varchar password
  }
  VENDEDOR {
    int id_vendedor PK,FK
    bigint telefono
  }
  ARTICULO {
    int id_articulo PK
    varchar titulo
    varchar descripcion
    condicion condicion
    categoria categoria
    text ruta_imagen
    numeric precio
    int stock
    timestamptz fecha_publicacion
    int id_vendedor FK
  }
  COMPRA {
    int id_compra PK
    int id_comprador FK
    timestamp fecha
    numeric total
  }
  DETALLE_COMPRA {
    int id_compra PK,FK
    int id_articulo PK,FK
    int cantidad
    numeric precio_historico_unitario
  }
  CARRITO {
    int id_usuario PK,FK
    int id_articulo PK,FK
    int cantidad
    timestamptz fecha_agregado
  }
  FAVORITO {
    int id_usuario PK,FK
    int id_articulo PK,FK
    timestamptz fecha
  }

  USUARIO ||--o| VENDEDOR : "es"
  VENDEDOR ||--o{ ARTICULO : "publica"
  USUARIO ||--o{ COMPRA : "realiza"
  COMPRA ||--|{ DETALLE_COMPRA : "contiene"
  ARTICULO ||--o{ DETALLE_COMPRA : "vendido en"
  USUARIO ||--o{ CARRITO : "tiene"
  ARTICULO ||--o{ CARRITO : "en"
  USUARIO ||--o{ FAVORITO : "guarda"
  ARTICULO ||--o{ FAVORITO : "en"
```
