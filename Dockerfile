FROM golang:1.24.2-alpine AS base

RUN apk add --no-cache \
    bash \
    curl \
    git \
    make \
    tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM base AS build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o payment-service-api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o payment-service-worker ./cmd/worker/main.go

FROM alpine:latest AS release

WORKDIR /app

COPY --from=build /app/config/container/start-app.sh ./
COPY --from=build /app/payment-service-api .
COPY --from=build /app/payment-service-worker .

EXPOSE 8080

CMD ["/app/start-app.sh"]
