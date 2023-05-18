FROM golang:1.20-alpine3.17 AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -o myapp ./cmd

FROM alpine:latest

COPY --from=builder /app/myapp .
COPY --from=builder /app/shared.kdbx .
COPY --from=builder /app/config.yaml .


CMD ["./myapp"]