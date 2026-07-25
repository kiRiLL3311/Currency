FROM golang:1.24

WORKDIR /app

COPY news/go.mod news/go.sum ./
RUN go mod download

COPY news/ .

RUN go build -o news-service .

EXPOSE 8080

CMD ["./news-service"]
