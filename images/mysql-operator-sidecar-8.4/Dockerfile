###############################################################################
#  Docker image for MySQL 8.4 sidecar containers
#
#  Unlike the 5.7/8.0 sidecar images (manually built bases, xtrabackup
#  compiled from source), this image installs the official Percona
#  XtraBackup 8.4 packages: the pxb-84-lts yum repo ships both x86_64 and
#  aarch64 for EL9, so no source build is needed for multi-arch.
#  xtrabackup 8.0 cannot back up MySQL 8.4 servers, hence the new image.
###############################################################################

# build mysql-operator-sidecar from the operator fork
FROM --platform=$TARGETPLATFORM golang:1.21 AS builder

ARG OPERATOR_REF=extra_image

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
# gzip is invoked by the sidecar backup/candidate streams
RUN dnf install -y https://repo.percona.com/yum/percona-release-latest.noarch.rpm \
    && percona-release enable pxb-84-lts release \
    && dnf install -y percona-xtrabackup-84 gzip unzip \
    && dnf clean all \
    && rm -rf /var/cache/dnf

# rclone via the official arch-aware installer
RUN curl -fsSL https://rclone.org/install.sh | bash \
    && rclone version

COPY --from=builder /mysql-operator-sidecar /usr/local/bin/mysql-operator-sidecar
RUN chmod 755 /usr/local/bin/mysql-operator-sidecar /usr/local/bin/docker-entrypoint.sh \
    && chown 999:999 /usr/local/bin/mysql-operator-sidecar

USER mysql
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
