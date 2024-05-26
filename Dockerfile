FROM golang:1.22.3-alpine

# Установка зависимостей и копирование исходного кода
WORKDIR /GraphQLOzon
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

# Установка PostgreSQL только если необходимо
RUN apk add --no-cache postgresql postgresql-contrib

# Копирование файла со скриптом в рут директорию и изменение прав доступа к нему
RUN chmod +x /GraphQLOzon/start.sh

# Стартовый скрипт
CMD ["/GraphQLOzon/start.sh"]
