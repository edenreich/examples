FROM golang:1.19.2-alpine3.16 AS common
WORKDIR /app

FROM common AS builder
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o email .

FROM common AS development
WORKDIR /app

FROM scratch AS production
COPY --from=builder /app/email /email
ENTRYPOINT [ "/email" ]
