FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrate.sh /clickhouse/migrate.sh
RUN chmod +x /clickhouse/migrate.sh
WORKDIR /clickhouse

CMD ["/clickhouse/migrate.sh"]