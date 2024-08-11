FROM golang:1.23rc2-alpine3.20

WORKDIR /register_service

COPY . .

RUN go mod tidy

RUN go build -o register_service register.go

CMD ["./register_service"]