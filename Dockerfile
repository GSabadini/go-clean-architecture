FROM golang:alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    MONGODB_HOST=mongodb \
    MONGODB_DATABASE=bank \
    POSTGRES_HOST=postgres \
    POSTGRES_PORT=5432 \
    POSTGRES_USER=dev \
    POSTGRES_PASSWORD=dev \
    POSTGRES_DATABASE=bank \
    POSTGRES_DRIVER=postgres \
    APP_NAME=go-bank-transfer \
    APP_PORT=3001

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -a --installsuffix cgo --ldflags="-s" -o main

FROM scratch

COPY --from=builder /build .

ENTRYPOINT ["./main"]

EXPOSE 3001