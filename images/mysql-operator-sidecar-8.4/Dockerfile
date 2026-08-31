###############################################################################
#  Docker image for MySQL 8.4 sidecar containers
#
#  Unlike the 5.7/8.0 sidecar images (manually built bases, xtrabackup
#  compiled from source), this image installs the official Percona
#  XtraBackup 8.4 packages: the pxb-84-lts yum repo ships both x86_64 and
#  aarch64 for EL9, so no source build is needed for multi-arch.
#  xtrabackup 8.0 cannot back up MySQL 8.4 servers, hence the new image.
###############################################################################

# build mysql-operator-sidecar from the operator fork; the Go version also
# pins the stdlib compiled into the binary - keep it current or the trivy
# CRITICAL gate rejects the image (the module itself builds fine on new Go)
FROM --platform=$TARGETPLATFORM golang:1.26-bookworm AS builder

# Release convention: tag the mysql-operator repo first (e.g. v0.8.0), then
# dispatch the build with build-args OPERATOR_REF=<that tag> so the sidecar
# binary is pinned to an exact commit (git clone --branch accepts tags).
# The default is pinned to the latest release tag; bump it per release.
ARG OPERATOR_REF=v0.8.0

# git is bundled with the golang image (wget/unzip are not)
RUN git clone --depth 1 --branch ${OPERATOR_REF} \
      https://github.com/ksmartdata/mysql-operator.git /build/mysql-operator
WORKDIR /build/mysql-operator
RUN CGO_ENABLED=0 go build -o /mysql-operator-sidecar ./cmd/mysql-operator-sidecar

# runtime: EL9 (matches Percona's own PXB 8.4 image base)
FROM --platform=$TARGETPLATFORM oraclelinux:9

# -o: uid/gid 999 must match the other sidecar images (the operator runs the
# pod with this identity), but EL9 already assigns 999 to a system account
RUN groupadd -o -g 999 mysql \
    && useradd -o -u 999 -r -g 999 -s /sbin/nologin -c "Default Application User" mysql

COPY rootfs/ /

# percona-xtrabackup-84 provides xtrabackup, xbstream and xbcloud;
# gzip is invoked by the sidecar backup/candidate streams;
# perl-DBI/perl-DBD-mysql are the runtime of the percona-toolkit scripts below
RUN dnf install -y https://repo.percona.com/yum/percona-release-latest.noarch.rpm \
    && percona-release enable pxb-84-lts release \
    && dnf install -y percona-xtrabackup-84 gzip unzip perl perl-DBI perl-DBD-mysql \
    && dnf clean all \
    && rm -rf /var/cache/dnf

# pt-heartbeat / pt-kill: the operator runs both as extra pod containers off
# this same sidecar image. Percona ships no aarch64 EL9 percona-toolkit rpm,
# but these two tools are pure Perl scripts - take them from the pinned
# release tarball instead of the arch-gated package
ARG PERCONA_TOOLKIT_VERSION=v3.7.1
RUN curl -fsSL https://github.com/percona/percona-toolkit/archive/refs/tags/${PERCONA_TOOLKIT_VERSION}.tar.gz \
      | tar -xz -C /tmp \
    && install -m 755 /tmp/percona-toolkit-*/bin/pt-heartbeat /tmp/percona-toolkit-*/bin/pt-kill /usr/local/bin/ \
    && rm -rf /tmp/percona-toolkit-* \
    && pt-heartbeat --version \
    && pt-kill --version

# rclone: pinned instead of the upstream installer, so the binary in the
# image is reproducible and the trivy CRITICAL gate has a version to name
ARG TARGETARCH
ARG RCLONE_RELEASE=1.74.3
RUN target_arch="${TARGETARCH:-$(uname -m)}" \
    && rclone_ver="${RCLONE_RELEASE#v}" \
    && case "${target_arch}" in \
        amd64|x86_64) rclone_arch=amd64 ;; \
        arm64|aarch64) rclone_arch=arm64 ;; \
        *) echo "Unsupported rclone architecture: ${target_arch}" >&2; exit 1 ;; \
    esac \
    && curl -fsSLo /tmp/rclone.zip "https://downloads.rclone.org/v${rclone_ver}/rclone-v${rclone_ver}-linux-${rclone_arch}.zip" \
    && unzip -q /tmp/rclone.zip -d /tmp \
    && install -m 755 "/tmp/rclone-v${rclone_ver}-linux-${rclone_arch}/rclone" /usr/local/bin/rclone \
    && rclone version \
    && rm -rf /tmp/rclone.zip "/tmp/rclone-v${rclone_ver}-linux-${rclone_arch}"

# re-declared so the builder-stage ARG is visible here; records which
# mysql-operator ref the embedded sidecar binary was built from
ARG OPERATOR_REF=v0.8.0
LABEL io.ksmartdata.operator-ref=${OPERATOR_REF}

COPY --from=builder /mysql-operator-sidecar /usr/local/bin/mysql-operator-sidecar
RUN chmod 755 /usr/local/bin/mysql-operator-sidecar /usr/local/bin/docker-entrypoint.sh \
    && chown 999:999 /usr/local/bin/mysql-operator-sidecar

USER mysql
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
