FROM golang:1.22.3-alpine

# Установка зависимостей и копирование исходного кода
WORKDIR /GraphQLOzon
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

CMD go run ./cmd/main.go
