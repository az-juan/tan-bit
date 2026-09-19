
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

[Documentación del proyecto](./docs/)

## Dependencias
+ atlas: `curl -sSf https://atlasgo.sh | sh`
+ docker
+ docker compose
+ figlet
+ go 1.26.5
+ gum
+ sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
+ templ: `go install github.com/a-h/templ/cmd/templ@latest`

Para ejecutar la aplicación:
```bash
  chmod u+x run.sh
  ./run.sh
```
Para acceder a la página: http://localhost:8080

## Asegurarse de
+ Estar en el grupo 'docker', en caso de no estarlo, ejecutar: `sudo usermod -aG docker $USER`
+ Tener los puertos 5432 y 8080 libres
