FROM golang:alpine AS builder
WORKDIR /build
COPY go.mod .
RUN go mod download -x
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.11
RUN swag init
RUN go build -o auth-ms main.go

FROM alpine
WORKDIR /build
COPY --from=builder build/auth-ms build/auth-ms
EXPOSE 8080
ENTRYPOINT ["build/auth-ms"]