FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0


WORKDIR /build/app/shared

COPY ../shared .

WORKDIR /build/app/src

COPY ./auth-ms/go.mod .
RUN go mod tidy
RUN go mod download -x
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.11
COPY ./auth-ms .
RUN swag init --parseDependency
RUN go build -o output main.go

FROM alpine
WORKDIR /build
COPY --from=builder build/app/src/output build/auth-ms
EXPOSE 8080
ENTRYPOINT ["build/auth-ms"]