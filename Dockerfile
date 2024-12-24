FROM golang:1.23-alpine AS builder

COPY . /github.com/merynayr/auth/source/
WORKDIR /github.com/merynayr/auth/source/

RUN go mod download
RUN go build -o ./bin/test_server cmd/grpc_server/main.go

FROM alpine:latest


WORKDIR /root/

COPY --from=builder /github.com/merynayr/auth/source/bin/test_server .

COPY local.env .

CMD ["./test_server", "-config-path=local.env"]