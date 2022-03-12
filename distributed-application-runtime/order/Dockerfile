FROM golang:1.17.3-alpine3.15 AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
COPY src src
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -ldflags="-w -s" -v -o bin/order \
    src/main.go

FROM golang:1.17.3-alpine3.15 AS development
WORKDIR /usr/src/app
RUN apk --update add curl \
    && curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
CMD [ "air" ]
EXPOSE 8080

FROM scratch AS production
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/bin/order ./bin/order
CMD [ "./bin/order" ]
EXPOSE 8080
