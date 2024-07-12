FROM golang:1.22 AS build

WORKDIR /build

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd


FROM alpine:latest

COPY --from=build /build/app ./app

ENTRYPOINT [ "./app" ]