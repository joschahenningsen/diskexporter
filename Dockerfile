FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o diskexporter main.go

FROM alpine:latest
COPY --from=builder /app/diskexporter /diskexporter

EXPOSE 1971

ENTRYPOINT ["/diskexporter"]
