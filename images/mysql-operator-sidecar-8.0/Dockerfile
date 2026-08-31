# keep the Go version current: the stdlib compiled into the sidecar binary
# is what the trivy CRITICAL gate scans (see the 8.4 image)
FROM --platform=$TARGETPLATFORM golang:1.26-bookworm as builder

USER root

RUN apt update -y

RUN apt install -y wget gcc unzip
ENV GOPROXY=https://goproxy.cn
ENV GOPATH=/build

# pinned mysql-operator ref (tag or commit - repo policy forbids pulling
# from branch heads); bump via build-args when cutting a new release
ARG OPERATOR_REF=v0.8.0

# build mysql-operator-sidecar
WORKDIR /build
RUN wget -O src.zip https://github.com/ksmartdata/mysql-operator/archive/${OPERATOR_REF}.zip
RUN unzip src.zip && mv mysql-operator-* mysql-operator-src
WORKDIR /build/mysql-operator-src
RUN go mod tidy
RUN CGO_ENABLED=0 go build  -o /mysql-operator-sidecar  cmd/mysql-operator-sidecar/main.go

# rclone: pinned instead of the upstream installer, so the binary in the
# image is reproducible and the trivy CRITICAL gate has a version to name
ARG TARGETARCH
ARG RCLONE_RELEASE=1.74.3
RUN target_arch="${TARGETARCH:-$(dpkg --print-architecture)}" \
    && rclone_ver="${RCLONE_RELEASE#v}" \
    && case "${target_arch}" in \
        amd64) rclone_arch=amd64 ;; \
        arm64) rclone_arch=arm64 ;; \
        *) echo "Unsupported rclone architecture: ${target_arch}" >&2; exit 1 ;; \
    esac \
    && curl -fsSLo /tmp/rclone.zip "https://downloads.rclone.org/v${rclone_ver}/rclone-v${rclone_ver}-linux-${rclone_arch}.zip" \
    && unzip -q /tmp/rclone.zip -d /tmp \
    && install -m 755 "/tmp/rclone-v${rclone_ver}-linux-${rclone_arch}/rclone" /usr/bin/rclone \
    && rclone version \
    && rm -rf /tmp/rclone.zip "/tmp/rclone-v${rclone_ver}-linux-${rclone_arch}"

FROM --platform=$TARGETPLATFORM ghcr.io/ksmartdata/mysql-operator-sidecar-8.0:v0.7.4
USER root
# the v0.7.4 base ships a full Go toolchain at /usr/local/go; a runtime
# image needs no compiler and trivy gates on its stale stdlib CVEs
RUN rm -rf /usr/local/go
COPY --from=builder /mysql-operator-sidecar /usr/local/bin/
COPY --from=builder /usr/bin/rclone /usr/local/bin/
RUN chmod +x /usr/local/bin/mysql-operator-sidecar
RUN chmod +x /usr/local/bin/rclone
RUN chown 999.999 /usr/local/bin/mysql-operator-sidecar

