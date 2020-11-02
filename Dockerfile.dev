FROM golang:1.14-stretch

WORKDIR /app

COPY . .

RUN go mod download && go get github.com/cespare/reflex

COPY _scripts/reflex/reflex.conf /

EXPOSE 3001

ENTRYPOINT ["reflex", "-c", "./_scripts/reflex/reflex.conf"]