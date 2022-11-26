FROM golang:1.19 as builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ../migrate .
RUN CGO_ENABLED=0 go build cmd/migrate/migrate.go

FROM scratch

WORKDIR /app

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/migrate .

ENV TZ UTC
ENTRYPOINT ["/app/migrate"]
