# stage 0
FROM golang:alpine as builder

WORKDIR /go/src/utarwyn.fr/storage-server
COPY . .

RUN mkdir ./bin && \
    apk add build-base upx && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s" -tags netgo -installsuffix netgo -o ./bin/storage && \
    upx -9 ./bin/storage

# stage 1
FROM scratch
WORKDIR /
COPY --from=builder /go/src/utarwyn.fr/storage-server/bin/ .
EXPOSE 8043
ENTRYPOINT ["/storage"]

