
# Mi Primera aplicacion web, TAN-BIT

![Alt text](./images/tanbit_logo.png)

Mercado de hardware/tecnología de segunda mano, destinado a facilitar la circulación de articulos que aún tienen utilidad.

Por cada artículo:
- Título
- Precio
- Descripción
- Condición
- Imagen/es
- Categoría
- Stock
- Contacto

## Dependencias
+ docker
+ go 1.26.5
+ sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
+ templ: `go install github.com/a-h/templ/cmd/templ@latest`

Para ejecutar la aplicación:
```bash
  make run
```
Una vez corriendo, si todavía no se realizó, aplicar migraciones:
```bash
  make status
  make apply
```
Luego, ejecutar tests:
```bash
  make test
```
Para detener la aplicación:
```bash
  make stop
```

## Asegurarse de:
+ Estar en el grupo 'docker', en caso de no estarlo, ejecutar: `sudo usermod -aG docker $USER`
+ Tener el puerto 5432 libre

Para acceder a la página: http://localhost:8080
