FROM --platform=$TARGETPLATFORM mongo:4.2.24 as all

USER root
ARG TARGETARCH
ARG RCLONE_RELEASE=1.74.3

COPY rootfs/ /

# To fix the error: The following signatures were invalid: EXPKEYSIG 4B7C549A058F8B6B MongoDB 4.2 Release Signing Key <packaging@mongodb.com>
RUN apt-get update --allow-insecure-repositories \
    && apt-get install -y unzip curl bash

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN target_arch="${TARGETARCH:-$(dpkg --print-architecture)}" \
    && rclone_ver="${RCLONE_RELEASE#v}" \
    && case "${target_arch}" in \
        amd64) rclone_arch=amd64 ;; \
        arm64) rclone_arch=arm64 ;; \
        *) echo "Unsupported rclone architecture: ${target_arch}" >&2; exit 1 ;; \
    esac \
    && echo "TARGETARCH=${TARGETARCH}, target_arch=${target_arch}, rclone_arch=${rclone_arch}, rclone_ver=${rclone_ver}" \
    && curl -fsSLo /tmp/rclone.zip "https://downloads.rclone.org/v${rclone_ver}/rclone-v${rclone_ver}-linux-${rclone_arch}.zip" \
    && unzip -q /tmp/rclone.zip -d /tmp \
    && cp "/tmp/rclone-v${rclone_ver}-linux-${rclone_arch}/rclone" /usr/local/bin/rclone \
    && chmod 755 /usr/local/bin/rclone \
    && test -x /usr/local/bin/rclone \
    && (/usr/local/bin/rclone version || true) \
    && rm -rf /tmp/rclone.zip "/tmp/rclone-v${rclone_ver}-linux-${rclone_arch}"

RUN cp /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime

RUN chmod +x /usr/local/bin/docker-entrypoint.sh
