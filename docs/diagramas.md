# Diagramas de arquitectura — TAN-BIT

```mermaid
flowchart TB
  subgraph app["tan-bit.com/tan-bit (Go module)"]
    main["main.go\n(entrypoint, routing, static fs)"]
    handlers["handlers/\nArticuloHandler, VendedorHandler,\nUsuarioHandler"]
    views["views/*_templ.go\n(Marketplace, Carrito, Compra,\nVender, Registrarse, MiCuenta..."]
    dbpkg["db/"]
    static["static/\ncss/*.css, assets/*"]
  end
  subgraph dbpkg2["db/"]
    conn["connection/Open()"]
    sqlc["queries (codegen)\nGetArticulos, CreateCompra, GetCarritoCompra,\nAddFavorito, GetComprasByUsuario, ..."]
    queries["queries/ *.sql"]
    schema["schema.sql"]
    migrations["migrations/ (Atlas)"]
  end
  main --> handlers
  main --> views
  main --> conn
  handlers --> views
  handlers --> sqlc
  conn --> sqlc
  queries -.sqlc generate.-> sqlc
  schema -.Atlas.-> migrations
```

## Componentes y conectores (C&C view, runtime)

```mermaid
flowchart LR
  browser["Browser\n(HTML/CSS)"]
  http["Go HTTP Server\n:8080 (http.ServeMux)"]
  handler["ArticuloHandler\nVendedorHandler\nUsuarioHandler"]
  renderer["Templ renderer\nviews/*_templ.go"]
  sqlc["sqlc queries\npgx pool"]
  pg["PostgreSQL\n:5432"]
  fs["Static file server\n/static/*"]
  uploads["Served dir\n/app/uploads\n(nginx)"]

  browser -- "HTTP GET/POST" --> http
  http --> handler
  http --> fs
  handler --> renderer
  handler --> sqlc
  sqlc -- "SQL (pgx)" --> pg
  renderer -- "HTML" --> browser
  fs -- "CSS/assets" --> browser
  uploads -- /uploads/* --> browser
```

## Diagrama de flujo de los handlers (flujo de petición)

```mermaid
flowchart TB
  client["Browser"]
  mux["http.ServeMux\n(de main.go)"]
  arth["ArticuloHandler\nServeHTTP"]
  vendh["VendedorHandler\nServeHTTP"]
  usuh["UsuarioHandler\nServeHTTP"]

  subgraph Market["GET /"]
    mkt["ArticuloHandler.MarketplaceHandler\nGetArticulos → views.Marketplace"]
  end
  subgraph Carrito["GET /carrito"]
    car["ArticuloHandler.CarritoHandler\nGetCarritoByUsuario → views.Carrito"]
  end
  subgraph Compra["GET/POST /compra"]
    comp["ArticuloHandler.CompraHandler\nGetCarritoCompra + CreateCompra + CreateDetalleCompra"]
  end
  subgraph Fav["GET /favoritos"]
    fav["ArticuloHandler.FavoritosHandler\nGetFavoritosByUsuario → views.Favoritos"]
  end
  subgraph Vender["POST /vender (vendedor)"]
    vend["VendedorHandler.VenderHandler\nvalida campos → upload → CreateArticulo"]
  end
  subgraph Login["GET/POST /login, /registrarse, /mi_cuenta"]
    usu["UsuarioHandler\nLogin, Registrarse, MiCuenta"]
  end

  client --> mux
  mux --> arth
  mux --> vendh
  mux --> usu
  arth --> mkt
  arth --> car
  arth --> comp
  arth --> fav
  vendh --> Vender
  usu --> Login

  mkt --> sqlc
  comp --> sqlc
  fav --> sqlc
  Vender --> sqlc
  Login --> sqlc
  sqlc --> pg
  templ --> browser
```

## Modelo de datos (ERD)

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
    text direccion_xmr
    text cbu
    text alias
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

---

## Notas de implementación

- **No hay JSON**. La capa de presentación es HTML generado con templ.
- **Autenticación**: cookie `user_id` (valor entero = `id_usuario`). El handler determina si el usuario es vendedor si el registro de `vendedor` existe.
- **Validación de formularios**: en `VenderHandler` y `RegistrarseHandler`; los errores se muestran como `errMsg` en la vista.
- **Subida de imágenes**: `sanitize()` + `os.MkdirAll("uploads", 0755)` + `io.Copy`. La ruta guardada es `/uploads/<nombre>`.
- **Migrations**: Atlas (`db/migrations/*.sql`). También los ficheros `.sql` en `db/schema/schema.sql` y `db/queries/*.sql` se usan para codegen de sqlc.
- **Docker**: `Dockerfile.prod` (builder + final alpine), `compose.prod.yaml` (api/db/nginx), `run-app.sh` (dev). Persistencia XMR: `db/uploads`, PG: `pgdata`.
