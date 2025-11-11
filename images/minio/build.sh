#!/usr/bin/env bash
set -euxo pipefail
bash build-mc.sh ${MINIO_VERSION}
bash build-minio.sh ${MINIO_VERSION}
cp Dockerfile.release /tmp/minio/Dockerfile.release
