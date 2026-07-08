BUILD_DIR="/tmp/build"
function _build() {
	local component=$1

	cd ${CLONE_DIR}
  LDFLAGS=$(go run buildscripts/gen-ldflags.go)
  mkdir -p "${BUILD_DIR}"
	CGO_ENABLED=0 go build -tags kqueue -trimpath --ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${component}
	chmod +x ${BUILD_DIR}/${component}
}

# $1: the minio tag
# $2: mc or minio
function main() {
  _clone ${1}
  _build ${2}
}