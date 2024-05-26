FROM golang:1.22.3

# Установка зависимостей и копирование исходного кода
WORKDIR /GraphQLOzon
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

# Установка PostgreSQL только если необходимо
RUN apk add --no-cache postgresql postgresql-contrib

# Стартовый скрипт
CMD ["sh", "start.sh"]
