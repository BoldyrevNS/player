FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o auth-ms main.go

FROM alpine

WORKDIR /build

COPY --from=builder build/auth-ms build/auth-ms

ENTRYPOINT ["build/auth-ms"]