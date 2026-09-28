FROM golang:1.26-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN go install github.com/air-verse/air@latest
RUN apk add curl && curl -sSf https://atlasgo.sh | sh
COPY . .
EXPOSE 8080
CMD ["air"]
