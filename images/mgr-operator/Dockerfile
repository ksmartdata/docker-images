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
RUN go mod download

COPY --from=downloader /tmp/mysql-operator-extra_image/images/mgr-diagnose/ /mgr-diagnose/
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-w -s" -o /usr/local/bin/diagnose daocloud.io/mcamel/mgr-diagnose

FROM --platform=$TARGETPLATFORM ghcr.io/ksmartdata/community-operator:${BASE_VERSION}
ARG BASE_VERSION
COPY --from=downloader /tmp/mgr-operator-package/mysqloperator/ /usr/lib/mysqlsh/python-packages/mysqloperator/
COPY --from=builder /usr/local/bin/diagnose /usr/local/bin/diagnose
# 1.35.4.post1是为来修复在报429的情况下operator不重试的bug，因为mgr-operator因为依赖特定版本的python没有办法用新版本的kopf
# 所以在特定版本的kopf上修改了kopf的代码来修复这个问题
USER root
RUN pip3 uninstall -y kopf && \
    pip3 install kopf-retry-after==1.35.4.post1
USER 2