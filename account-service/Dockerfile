FROM golang:1.18

WORKDIR /app

COPY . .

RUN go build -o account-service main.go

EXPOSE 8081

CMD ["./account-service"]
