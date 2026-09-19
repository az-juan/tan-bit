# Dockerfile multi-stage (evita dejar herramientas de compilación pesadas en la imagen final)

## Etapa 1: compilación

##Descarga la imagen base para usarla durante la construccion
FROM golang:1.26-alpine AS builder
##Establece el directorio de trabajo dentro del contenedor, y se coloca ahi
WORKDIR /app
##Copia las dependencias antes que el codigo en ./ (/app)
COPY go.mod go.sum ./
##Descarga modulos de go requeridos
RUN go mod download
##Copia el codigo fuente del proyecto en el contenedor (/api /app)
COPY . .
# Compila el code en un binario estático independiente y lo guarda en /app/api
RUN CGO_ENABLED=0 go build -o /app/api .

## Etapa 2: imagen final (solo el binario)

##Inicia desde una imagen minima, solo tiene el binario en /app/api
FROM alpine:3.24
##Directorio de ejecucion dentro del contenedor final
WORKDIR /app
##Crea usuario sin privilegios
RUN adduser -D -u 1000 user
##Copia el binario compilado en la etapa 1 a /usr/local/bin/api
COPY --from=builder /app/api /usr/local/bin/api
##Copia la carpeta en /app/static
COPY static /app/static
##Ejecuta sin privilegios
USER user
##Expone puerto
EXPOSE 8080
##Comando a ejecutar al iniciar contenedor (:17 lo puso en PATH)
CMD ["api"]
