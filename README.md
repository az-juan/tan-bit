
# Mi Primera aplicacion web, TAN-BIT

![Alt text](./images/tanbit_logo.png)

Mercado de hardware/tecnología de segunda mano, destinado a facilitar la circulación de articulos que aún tienen utilidad entre vecinos.

DESCENTRALIZADO - CONFIANZA - 

[Documentación del proyecto](./docs/)

## Dependencias
+ docker
+ docker compose
+ go
+ sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
+ templ: `go install github.com/a-h/templ/cmd/templ@latest`

Para ejecutar la aplicación:
```bash
  chmod u+x test-api.sh
  ./test-api.sh
```
Para acceder a la página: http://localhost:8080

## Asegurarse de
+ Estar en el grupo 'docker', en caso de no estarlo, ejecutar: `sudo usermod -aG docker $USER`.
Al ejecutar `groups` tiene que listarse "docker", de lo contrario cierra la sesión actual e inicia otra.
+ Tener puertos 5433 y 8080 libres
+ Iniciar el Docker Daemon con: `sudo systemctl start docker`
