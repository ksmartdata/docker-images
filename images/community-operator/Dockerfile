ARG MGR_OPERATOR_REF=8.4.9-2.1.11-r1
FROM alpine:3.22.2 AS downloader

ARG MGR_OPERATOR_REF
WORKDIR /tmp
RUN wget -q "https://github.com/ksmartdata/mgr-operator/archive/${MGR_OPERATOR_REF}.zip" -O mgr-operator.zip && \
    unzip -q mgr-operator.zip && \
    mv mgr-operator-* mgr-operator-src

FROM --platform=$TARGETPLATFORM python:3.13-slim AS python-deps

WORKDIR /tmp
COPY --from=downloader /tmp/mgr-operator-src/docker-deps/requirements.txt .
RUN python -m pip install --no-cache-dir --target=/tmp/site-packages -r requirements.txt --no-deps

FROM --platform=$TARGETPLATFORM container-registry.oracle.com/os/oraclelinux:9-slim

ARG MYSQL_REPO_URL=https://repo.mysql.com
ARG MYSQL_CONFIG_PKG=mysql84-community-release
ARG MYSQL_SHELL_REPO=mysql-tools-8.4-lts-community
ARG MYSQL_SHELL_VERSION=8.4.9

RUN rpm -U "${MYSQL_REPO_URL}/${MYSQL_CONFIG_PKG}-el9.rpm" && \
    printf "[main]\nip_resolve=4\ntimeout=30\nretries=3\ntsflags=nodocs\n" > /etc/dnf/dnf.conf && \
    microdnf install --disablerepo="*" --enablerepo=ol9_baseos_latest --enablerepo=ol9_appstream --enablerepo="${MYSQL_SHELL_REPO}" -y glibc-langpack-en "mysql-shell-${MYSQL_SHELL_VERSION}" && \
    rpm -e "${MYSQL_CONFIG_PKG}" && \
    microdnf clean all

RUN groupadd -g27 mysql && useradd -u27 -g27 mysql

RUN mkdir /mysqlsh && chown 2 /mysqlsh

COPY --from=python-deps /tmp/site-packages/ /usr/lib/mysqlsh/python-packages/
COPY --from=downloader /tmp/mgr-operator-src/mysqloperator/ /usr/lib/mysqlsh/python-packages/mysqloperator/

RUN sed -i "s/available_replicas=None,/available_replicas=0,/" /usr/lib/mysqlsh/python-packages/kubernetes/client/models/v1_stateful_set_status.py

USER 2

ENV HOME=/mysqlsh
