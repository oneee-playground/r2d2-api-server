FROM golang:1.22 AS build

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd


FROM alpine:latest

COPY --from=build /build/app ./app

ENTRYPOINT [ "./app" ]