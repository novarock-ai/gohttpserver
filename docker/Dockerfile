FROM golang:1.18 AS builder
WORKDIR /app/gohttpserver
ADD . /app/gohttpserver
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-X main.VERSION=docker' -o gohttpserver

FROM alpine

WORKDIR /app
COPY common/assets .
COPY --from=builder /app/gohttpserver/gohttpserver .
EXPOSE 9999
ENTRYPOINT [ "/bin/sh", "-c", "--", "./gohttpserver -r $FILE_ROOT --port 9999 --pin-root --not-exist-auto-mkdir --custom-cdn https://static01.dp.tech/openfiles.mlops.dp.tech/v2/static"]
