FROM --platform=$TARGETPLATFORM golang:1.19.2 as builder

USER root

RUN apt update -y && \
    apt install -y wget gcc unzip

ENV GOPROXY=https://goproxy.cn
ENV GOPATH=/build

# build mysql-operator-sidecar
WORKDIR /build
RUN wget -O extra_image.zip https://github.com/ksmartdata/mysql-operator/archive/refs/heads/extra_image.zip
RUN unzip extra_image.zip
WORKDIR /build/mysql-operator-extra_image
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /mysql-operator-sidecar cmd/mysql-operator-sidecar/main.go

# install rclone
RUN curl https://rclone.org/install.sh | bash

FROM --platform=$TARGETPLATFORM jiuker/mysql-operator-sidecar-5.7:v0.6.2
USER root
COPY --from=builder /mysql-operator-sidecar /usr/local/bin/
COPY --from=builder /usr/bin/rclone /usr/local/bin/
# debug
RUN rclone --version
RUN chmod +x /usr/local/bin/mysql-operator-sidecar
RUN chown 999.999 /usr/local/bin/mysql-operator-sidecar

