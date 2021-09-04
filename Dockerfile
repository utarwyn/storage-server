# stage 0
FROM golang:alpine as builder

WORKDIR /go/src/companyofcube.fr/storage
COPY . .

RUN mkdir ./bin && \
    apk add build-base upx && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s" -tags netgo -installsuffix netgo -o ./bin/storage && \

    mkdir ./bin/etc && \
    ID=$(shuf -i 100-9999 -n 1) && \
    upx -9 ./bin/storage && \
    echo $ID && \
    echo "appuser:x:$ID:$ID::/sbin/nologin:/bin/false" > ./bin/etc/passwd && \
    echo "appgroup:x:$ID:appuser" > ./bin/etc/group

# stage 1
FROM scratch
WORKDIR /
COPY --from=builder /go/src/companyofcube.fr/storage/bin/ .
USER appuser
EXPOSE 8043
ENTRYPOINT ["/storage"]

