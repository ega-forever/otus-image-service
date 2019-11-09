FROM golang:1.13 as builder

WORKDIR /app

COPY cmd .
COPY internal .
COPY go.mod .
COPY main.go .

RUN CGO_ENABLED=0 && go build -o service

###

FROM ubuntu:18.04 as runner

WORKDIR /app
ENV REST_PORT=8080

COPY --from=builder /app/service .

EXPOSE 8080
CMD ["./service", "grpc_server", "database"]