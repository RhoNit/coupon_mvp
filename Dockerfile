FROM golang:1.21-alpine

WORKDIR /app

# install build dependencies
RUN apk add --no-cache gcc musl-dev

# install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# generate swagger docs (swag init)
RUN swag init -g main.go -o docs

RUN go build -o main .

EXPOSE 8080

CMD ["./main"] 