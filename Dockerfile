FROM golang:1.14.0-alpine3.11 as builder

RUN apk add --no-cache --virtual .build-deps \
    bash \
    gcc \
    git \
    musl-dev

RUN mkdir /go-bank-transfer
COPY . /go-bank-transfer
WORKDIR /go-bank-transfer

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
RUN adduser -S -D -H -h /go-bank-transfer main
USER main

FROM scratch

COPY --from=builder /go-bank-transfer /go-bank-transfer

WORKDIR /go-bank-transfer

EXPOSE 3001

CMD ["./main"]