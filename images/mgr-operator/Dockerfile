ARG BASE_VERSION=8.3.0-2.1.2
FROM alpine:3.22.2 AS builder
WORKDIR /tmp

RUN apk add --no-cache wget unzip && \
    wget https://github.com/ksmartdata/mgr-operator/archive/refs/heads/package.zip -O mgr-operator.zip && \
    unzip mgr-operator.zip

RUN ls /tmp/mgr-operator-package

FROM container-registry.oracle.com/mysql/community-operator:${BASE_VERSION}
ARG BASE_VERSION
COPY --from=builder /tmp/mgr-operator-package/mysqloperator/ /usr/lib/mysqlsh/python-packages/mysqloperator/