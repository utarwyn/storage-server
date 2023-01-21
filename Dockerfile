# stage 0
FROM --platform=$BUILDPLATFORM golang:alpine as builder

ARG TARGETPLATFORM

WORKDIR /go/src/github.com/utarwyn/storage-server
COPY . .

RUN mkdir ./bin && \
    apk add build-base upx && \
    GOOS=$(echo $TARGETPLATFORM | cut -f1 -d/) && \
    GOARCH=$(echo $TARGETPLATFORM | cut -f2 -d/) && \
    GOARM=$(echo $TARGETPLATFORM | cut -f3 -d/ | sed "s/v//" ) && \
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build ${BUILD_ARGS} -ldflags="-s" -tags netgo -installsuffix netgo -o ./bin/storage-server && \
    upx -9 ./bin/storage-server

# stage 1
FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/utarwyn/storage-server/bin/ .
EXPOSE 8043
ENTRYPOINT ["/storage-server"]
