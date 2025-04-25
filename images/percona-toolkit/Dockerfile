FROM --platform=$TARGETPLATFORM centos:8.3.2011

USER root

RUN cd /etc/yum.repos.d/ && \
    sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-* && \
    sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-* && \
    yum makecache

RUN yum update -y

RUN yum -y install unzip wget cmake openssl-devel libaio libaio-devel automake autoconf \
    bison libtool ncurses-devel libgcrypt-devel libev-devel libcurl-devel zlib-devel \
    vim-common gcc-toolset-11-gcc gcc-toolset-11-gcc-c++ gcc-toolset-11-binutils \
    cyrus-sasl-devel libudev-devel \
    perl-DBI perl-DBD-MySQL

RUN wget -O libkmip.zip https://github.com/Percona-Lab/libkmip/archive/0ecda33598838b67bb4bb7a0005c92eea8b7405a.zip
RUN unzip libkmip.zip
RUN cp -r libkmip-0ecda33598838b67bb4bb7a0005c92eea8b7405a/. extra/libkmip/

# 1. install go
ENV GO_VERSION=1.22.4
ENV GO_DOWNLOAD_URL=https://golang.org/dl/
ENV GO_TARBALL=go${GO_VERSION}.linux-amd64.tar.gz
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then \
        GO_TARBALL=go${GO_VERSION}.linux-amd64.tar.gz; \
    elif [ "$ARCH" = "armv7l" ]; then \
        GO_TARBALL=go${GO_VERSION}.linux-armv6l.tar.gz; \
    elif [ "$ARCH" = "aarch64" ]; then \
        GO_TARBALL=go${GO_VERSION}.linux-arm64.tar.gz; \
    else \
        echo "Unsupported architecture: $ARCH"; exit 1; \
    fi && \
    wget ${GO_DOWNLOAD_URL}${GO_TARBALL} && \
    tar -C /usr/local -xzf ${GO_TARBALL} && \
    rm ${GO_TARBALL}
ENV PATH=$PATH:/usr/local/go/bin
RUN go version

# 2. install percona-toolkit
RUN wget -O percona-toolkit.zip https://github.com/percona/percona-toolkit/archive/refs/heads/3.x.zip
RUN unzip percona-toolkit.zip
WORKDIR /build/percona-toolkit-3.x
RUN yum install perl-ExtUtils-CBuilder perl-ExtUtils-MakeMaker -y
RUN perl Makefile.PL
RUN make
RUN make install