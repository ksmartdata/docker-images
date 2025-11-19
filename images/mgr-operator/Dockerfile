ARG BASE_VERSION=8.3.0-2.1.2
FROM alpine:3.22.2 AS downloader

WORKDIR /tmp
RUN wget https://github.com/ksmartdata/mgr-operator/archive/refs/heads/package.zip -O mgr-operator.zip && \
    unzip mgr-operator.zip \

RUN wget https://github.com/ksmartdata/mysql-operator/archive/refs/heads/extra_image.zip -O extra_image.zip && \
    unzip extra_image.zip
RUN ls /tmp/mgr-operator-package
RUN ls /tmp/mysql-operator-extra_image

FROM golang:1.23.3 AS builder
COPY --from=downloader /tmp/mysql-operator-extra_image/images/mgr-diagnose/ /mgr-diagnose/
WORKDIR /mgr-diagnose
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /usr/local/bin/diagnose daocloud.io/mcamel/mgr-diagnose


FROM container-registry.oracle.com/mysql/community-operator:${BASE_VERSION}
ARG BASE_VERSION
COPY --from=downloader /tmp/mgr-operator-package/mysqloperator/ /usr/lib/mysqlsh/python-packages/mysqloperator/
COPY --from=builder /usr/local/bin/diagnose /usr/local/bin/diagnose
RUN chmod +x /usr/local/bin/diagnose