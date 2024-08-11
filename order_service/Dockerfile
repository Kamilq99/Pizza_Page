FROM golang:1.23rc2-alpine3.20

WORKDIR /order_service

COPY . .

RUN go mod tidy

RUN go build -o order_service order.go

CMD ["./order_service"]