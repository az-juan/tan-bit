##Descarga la imagen base para usarla durante la construccion
FROM golang:1.26-alpine
##Necesitamos make
RUN apk add --no-cache make
##air y sqlc si queremos que el contenedor las tenga
RUN go install github.com/air-verse/air@latest && \
        go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
##Establece el directorio de trabajo dentro del contenedor, y se coloca ahi
WORKDIR /app
##Copia las dependencias antes que el codigo en ./ (/app)
COPY go.mod go.sum ./
##Descarga modulos de go requeridos
RUN go mod download
##Copia el codigo fuente del proyecto en el contenedor (/api /app). Image build time
COPY . .
EXPOSE 8080
CMD ["make", "run"]
