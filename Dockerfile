FROM golang:1.12-stretch

WORKDIR /go-stone

COPY . .

RUN go mod download
RUN go get github.com/cespare/reflex

COPY reflex.conf /

EXPOSE 3001

ENTRYPOINT ["reflex", "-c", "/reflex.conf"]