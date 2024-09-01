FROM golang:1.18

WORKDIR /app

COPY . .

RUN go build -o order-service main.go

EXPOSE 8082

CMD ["./order-service"]
