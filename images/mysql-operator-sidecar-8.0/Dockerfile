FROM --platform=$TARGETPLATFORM golang:1.21.11 as builder

USER root

RUN apt update -y

RUN apt install -y wget gcc unzip
ENV GOPROXY=https://goproxy.cn
ENV GOPATH=/build

# build mysql-operator-sidecar
WORKDIR /build
RUN wget -O arm64.zip https://github.com/ksmartdata/mysql-operator/archive/refs/heads/extra_image.zip
RUN unzip arm64.zip
WORKDIR /build/mysql-operator-extra_image
RUN go mod tidy
RUN CGO_ENABLED=0 go build  -o /mysql-operator-sidecar  cmd/mysql-operator-sidecar/main.go

# install rclone
RUN curl https://rclone.org/install.sh | bash

FROM --platform=$TARGETPLATFORM ghcr.io/ksmartdata/mysql-operator-sidecar-8.0:v0.7.4
USER root
COPY --from=builder /mysql-operator-sidecar /usr/local/bin/
COPY --from=builder /usr/bin/rclone /usr/local/bin/
RUN chmod +x /usr/local/bin/mysql-operator-sidecar
RUN chmod +x /usr/local/bin/rclone
RUN chown 999.999 /usr/local/bin/mysql-operator-sidecar

