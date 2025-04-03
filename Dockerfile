FROM golang:1.24 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./

RUN go mod tidy

COPY ../ ./

RUN go build -o gpsd-map-mgmt ./cmd/main.go

FROM debian:bookworm

WORKDIR /app

RUN apt-get update && apt-get install -y curl lsof net-tools ca-certificates libc6 iputils-ping

COPY --from=builder /app/gpsd-map-mgmt .

COPY ../config/ config/

EXPOSE 7000

CMD ["./gpsd-map-mgmt"]
