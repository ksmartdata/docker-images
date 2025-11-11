#!/usr/bin/env bash
set -euxo pipefail

source ./base.sh

COMPONENT='mc'
CLONE_DIR="/tmp/${COMPONENT}"
RELEASE_VERSION_FILE="/tmp/${COMPONENT}.txt"
export MC_RELEASE='RELEASE'

function _clone() {
   if [[ $# -ne 1 ]]; then
       echo "Usage: $0 <minio release> (e.g.  RELEASE.2025-09-07T16-13-09Z)"
       exit 1
   fi
   go run . -minio_release=${1}
   git clone --branch $(cat ${RELEASE_VERSION_FILE}) --depth 1 https://github.com/minio/${COMPONENT}.git ${CLONE_DIR}
}

main "$@" ${COMPONENT}


