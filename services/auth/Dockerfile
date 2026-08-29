FROM golang:1.26-alpine3.23 AS builder

WORKDIR /auth

COPY . . /auth/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/

FROM alpine3.23

WORKDIR /app

COPY --from=builder /bin/app /app/bin

ENTRYPOINT [ "./app/bin" ]