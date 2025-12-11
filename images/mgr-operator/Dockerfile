ARG BASE_VERSION=8.3.0-2.1.2
FROM alpine:3.22.2 AS downloader

WORKDIR /tmp
RUN wget -q https://github.com/ksmartdata/mgr-operator/archive/refs/heads/package.zip -O mgr-operator.zip && \
    unzip -q mgr-operator.zip && \
    wget -q https://github.com/ksmartdata/mysql-operator/archive/refs/heads/extra_image.zip -O extra_image.zip && \
    unzip -q extra_image.zip

FROM --platform=$TARGETPLATFORM golang:1.23.3-alpine AS builder
USER root
WORKDIR /mgr-diagnose
COPY --from=downloader /tmp/mysql-operator-extra_image/images/mgr-diagnose/go.mod ./
COPY --from=downloader /tmp/mysql-operator-extra_image/images/mgr-diagnose/go.sum ./
RUN go mod tidy && go mod download

COPY --from=downloader /tmp/mysql-operator-extra_image/images/mgr-diagnose/ /mgr-diagnose/
RUN go mod tidy
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-w -s" -o /usr/local/bin/diagnose daocloud.io/mcamel/mgr-diagnose

# Set base image according to platform architecture
FROM --platform=$TARGETPLATFORM ghcr.io/ksmartdata/community-operator:${BASE_VERSION}
ARG BASE_VERSION
COPY --from=downloader /tmp/mgr-operator-package/mysqloperator/ /usr/lib/mysqlsh/python-packages/mysqloperator/
COPY --from=builder /usr/local/bin/diagnose /usr/local/bin/diagnose