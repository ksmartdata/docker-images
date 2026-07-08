FROM --platform=$TARGETPLATFORM  ubuntu:22.04
LABEL maintainer="Team ACID @ Zalando <team-acid@zalando.de>"

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ARG TARGETARCH
ARG RCLONE_RELEASE=1.74.3
RUN apt-get update     \
    && apt-get install --no-install-recommends -y \
        apt-utils \
        ca-certificates \
        lsb-release \
        pigz \
        python3-pip \
        python3-setuptools \
        curl \
        jq \
        gnupg \
        gcc \
        libffi-dev \
        unzip \
    && target_arch="${TARGETARCH:-$(dpkg --print-architecture)}" \
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
    && rm -rf /tmp/rclone.zip "/tmp/rclone-v${rclone_ver}-linux-${rclone_arch}" \
    && curl -sL https://aka.ms/InstallAzureCLIDeb | bash \
    && pip3 install --upgrade pip \
    && pip3 install --no-cache-dir awscli --upgrade \
    && echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list \
    && cat /etc/apt/sources.list.d/pgdg.list \
    && curl --silent https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
    && apt-get update \
    && apt-get install --no-install-recommends -y  \
        postgresql-client-15  \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY * ./
RUN chmod +x ./actions.sh ./delete.sh ./dump.sh ./restore.sh

ENV PG_DIR=/usr/lib/postgresql

ENTRYPOINT ["/actions.sh"]
