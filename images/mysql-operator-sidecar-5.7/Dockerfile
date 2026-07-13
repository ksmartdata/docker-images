FROM --platform=$TARGETPLATFORM golang:1.21.11 as builder

USER root

RUN apt update -y && \
    apt install -y wget gcc unzip

ENV GOPROXY=https://goproxy.cn
ENV GOPATH=/build

# pinned mysql-operator ref (tag or commit - repo policy forbids pulling
# from branch heads); bump via build-args when cutting a new release
ARG OPERATOR_REF=v0.8.0

# build mysql-operator-sidecar
WORKDIR /build
RUN wget -O src.zip https://github.com/ksmartdata/mysql-operator/archive/${OPERATOR_REF}.zip
RUN unzip src.zip && mv mysql-operator-* mysql-operator-src
WORKDIR /build/mysql-operator-src
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /mysql-operator-sidecar cmd/mysql-operator-sidecar/main.go

# install rclone
RUN curl https://rclone.org/install.sh | bash

FROM --platform=$TARGETPLATFORM jiuker/mysql-operator-sidecar-5.7:v0.6.2
USER root
COPY --from=builder /mysql-operator-sidecar /usr/local/bin/
COPY --from=builder /usr/bin/rclone /usr/local/bin/
RUN chmod +x /usr/local/bin/mysql-operator-sidecar
RUN chmod +x /usr/local/bin/rclone
RUN chown 999.999 /usr/local/bin/mysql-operator-sidecar

