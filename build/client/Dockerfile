FROM golang:1.21-alpine3.18 as builder

WORKDIR /client

COPY . .

RUN apk add git

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/client/main.go

CMD ["./main"]