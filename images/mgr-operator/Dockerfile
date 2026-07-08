ARG BASE_VERSION=8.3.0-2.1.2
ARG BASE_IMAGE_REPOSITORY=ghcr.io/ksmartdata/community-operator
FROM alpine:3.22.2 AS downloader

ARG MGR_OPERATOR_REF=package
ARG MYSQL_OPERATOR_EXTRA_IMAGE_REF=extra_image
WORKDIR /tmp
RUN wget -q "https://github.com/ksmartdata/mgr-operator/archive/${MGR_OPERATOR_REF}.zip" -O mgr-operator.zip && \
    unzip -q mgr-operator.zip && \
    mv mgr-operator-* mgr-operator-src && \
    wget -q "https://github.com/ksmartdata/mysql-operator/archive/${MYSQL_OPERATOR_EXTRA_IMAGE_REF}.zip" -O extra_image.zip && \
    unzip -q extra_image.zip && \
    mv mysql-operator-* mysql-operator-src

FROM --platform=$TARGETPLATFORM golang:1.23.3-alpine AS builder
USER root
WORKDIR /mgr-diagnose
COPY --from=downloader /tmp/mysql-operator-src/images/mgr-diagnose/go.mod ./
COPY --from=downloader /tmp/mysql-operator-src/images/mgr-diagnose/go.sum ./
RUN go mod download

COPY --from=downloader /tmp/mysql-operator-src/images/mgr-diagnose/ /mgr-diagnose/
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-w -s" -o /usr/local/bin/diagnose daocloud.io/mcamel/mgr-diagnose

FROM --platform=$TARGETPLATFORM ${BASE_IMAGE_REPOSITORY}:${BASE_VERSION}
ARG BASE_VERSION
ARG BASE_IMAGE_REPOSITORY
COPY --from=downloader /tmp/mgr-operator-src/mysqloperator/ /usr/lib/mysqlsh/python-packages/mysqloperator/
COPY --from=builder /usr/local/bin/diagnose /usr/local/bin/diagnose
USER 2
