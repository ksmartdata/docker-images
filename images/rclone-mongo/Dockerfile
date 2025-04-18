FROM --platform=$TARGETPLATFORM mongo:4.2.24 as all

USER root

COPY rootfs/ /

# To fix the error: The following signatures were invalid: EXPKEYSIG 4B7C549A058F8B6B MongoDB 4.2 Release Signing Key <packaging@mongodb.com>
RUN apt-get update --allow-insecure-repositories

RUN apt-get install -y unzip curl bash

# Install latest gosu version to fix security vulnerabilities
RUN set -eux; \
    apt-get install -y ca-certificates wget; \
    GOSU_VERSION=1.17; \
    wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$(dpkg --print-architecture)"; \
    wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$(dpkg --print-architecture).asc"; \
    chmod +x /usr/local/bin/gosu; \
    gosu --version; \
    gosu nobody true

RUN curl https://rclone.org/install.sh | bash

RUN cp /usr/bin/rclone /usr/local/bin/

RUN cp /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime

RUN chmod +x /usr/local/bin/rclone

RUN chmod +x /usr/local/bin/docker-entrypoint.sh