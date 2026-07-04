FROM golang:1.25

WORKDIR /app

COPY auth/go.mod auth/go.sum ./
RUN go mod download

COPY auth/ .

RUN go build -o auth-service .

EXPOSE 8080

CMD ["./auth-service"]