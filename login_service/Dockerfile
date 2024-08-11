FROM golang:1.23rc2-alpine3.20

WORKDIR /login_service

COPY . .

RUN go mod tidy

RUN go build -o login_service login.go

CMD ["./login_service"]