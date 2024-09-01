FROM golang:1.18

WORKDIR /app

COPY . .

RUN go build -o coordinator-service main.go

EXPOSE 8080

CMD ["./coordinator-service"]
