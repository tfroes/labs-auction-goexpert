FROM golang:1.23.11-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/auction cmd/auction/main.go

# FROM scratch
FROM alpine:latest

COPY --from=builder /app/auction .
COPY --from=builder /app/cmd/auction/.env .env

EXPOSE 8080

ENTRYPOINT ["/auction"]