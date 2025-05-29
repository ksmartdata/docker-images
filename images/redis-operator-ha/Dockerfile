FROM --platform=$BUILDPLATFORM golang:1.23.9-alpine3.21 AS build
RUN apk --no-cache add \
  bash \
  unzip

WORKDIR /src
RUN wget https://github.com/ksmartdata/redis-operator-ha/archive/refs/heads/package.zip -O redis-operator-ha.zip
RUN unzip redis-operator-ha.zip

ARG TARGETOS 
ARG TARGETARCH 
ARG VERSION
WORKDIR /src/redis-operator-ha-package
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH VERSION=$VERSION ./scripts/build.sh

FROM alpine:latest
RUN apk --no-cache add \
  ca-certificates
COPY --from=build /src/redis-operator-ha-package/bin/redis-operator /usr/local/bin
RUN addgroup -g 1000 rf && \
  adduser -D -u 1000 -G rf rf && \
  chown rf:rf /usr/local/bin/redis-operator
USER rf

ENTRYPOINT ["/usr/local/bin/redis-operator"]