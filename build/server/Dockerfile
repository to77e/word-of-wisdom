FROM golang:1.21-alpine3.18 as builder

WORKDIR /app

COPY . .

RUN apk add git build-base binutils-gold

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

CMD ["./main"]


