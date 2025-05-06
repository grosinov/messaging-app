FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/server.go

EXPOSE 8080

ENTRYPOINT ["./app"]