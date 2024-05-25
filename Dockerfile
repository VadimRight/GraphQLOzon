FROM golang:1.22.1

WORKDIR /GraphQLOzon

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

CMD go run cmd/url-shortener/main.go
