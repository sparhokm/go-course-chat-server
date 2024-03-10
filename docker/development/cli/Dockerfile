FROM golang:1.22-alpine3.19

RUN apk update && \
    apk upgrade && \
    apk add git && \
    apk add bash && \
    apk add protobuf-dev && \
    rm -rf /var/cache/apk/*

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
RUN go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
RUN go install -mod=mod github.com/vektra/mockery/v2@v2.42.0
RUN go install github.com/envoyproxy/protoc-gen-validate@v1.0.4
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1

COPY ./common/wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

WORKDIR /app