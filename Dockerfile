FROM golang:1.21.0

WORKDIR /fitness-api

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build cmd/main.go

CMD ["./main"]