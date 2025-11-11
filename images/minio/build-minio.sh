#!/usr/bin/env bash
set -euxo pipefail

source ./base.sh

COMPONENT='minio'
CLONE_DIR="/tmp/${COMPONENT}"
export MINIO_RELEASE='RELEASE'

function _clone() {
   if [[ $# -ne 1 ]]; then
       echo "Usage: $0 <minio release> (e.g.  RELEASE.2025-09-07T16-13-09Z)"
       exit 1
   fi
   git clone --branch ${1} --depth 1 https://github.com/minio/${COMPONENT}.git ${CLONE_DIR}
}

main "$@" ${COMPONENT}