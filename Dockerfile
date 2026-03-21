FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .

RUN CGO_ENABLED=0 go build -mod=vendor -o /out/stock-app ./cmd

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/stock-app /app/stock-app
COPY templates /app/templates

EXPOSE 8080

USER app

CMD ["/app/stock-app"]
