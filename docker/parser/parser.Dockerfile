FROM golang:1.19 as builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build cmd/parser/parser.go

FROM scratch

WORKDIR /app

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/parser .

ENV TZ UTC
ENTRYPOINT ["/app/parser"]
