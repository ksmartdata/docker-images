#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

FROM --platform=$TARGETPLATFORM bellsoft/liberica-openjdk-alpine-musl:8

RUN apk add --no-cache bash gettext nmap-ncat openssl busybox-extras tzdata
# set default timezone
ENV TZ=Asia/Shanghai
RUN ln -fs /usr/share/zoneinfo/$TZ /etc/localtime

ARG version=5.3.3

# Rocketmq version
ENV ROCKETMQ_VERSION=${version}

# Rocketmq home
ENV ROCKETMQ_HOME=/root/rocketmq/nameserver

WORKDIR  ${ROCKETMQ_HOME}

# 4.3.2 的签名 key（duhengforever@apache.org，fpr B1BD9DB653F1B79C742839EA9087A12DD1BC4712）
# 从未收录进官方 KEYS 文件，这里以仓库内文件补充导入；发行包已另经 archive 的 md5/sha1 双重核对
COPY KEYS.duheng-4x /tmp/KEYS.duheng-4x

# Install
RUN set -eux; \
    apk add --virtual .build-deps curl gnupg unzip; \
    curl https://archive.apache.org/dist/rocketmq/${ROCKETMQ_VERSION}/rocketmq-all-${ROCKETMQ_VERSION}-bin-release.zip -o rocketmq.zip; \
    curl https://archive.apache.org/dist/rocketmq/${ROCKETMQ_VERSION}/rocketmq-all-${ROCKETMQ_VERSION}-bin-release.zip.asc -o rocketmq.zip.asc; \
	curl -L https://www.apache.org/dist/rocketmq/KEYS -o KEYS; \
	gpg --import KEYS; \
	gpg --import /tmp/KEYS.duheng-4x; \
    gpg --batch --verify rocketmq.zip.asc rocketmq.zip; \
    unzip rocketmq.zip; \
	mv rocketmq-all*/* . ; \
	rmdir rocketmq-all* ; \
	rm rocketmq.zip ; \
    rm rocketmq.zip.asc KEYS; \
	apk del .build-deps ; \
    rm -rf /var/cache/apk/* ; \
    rm -rf /tmp/*

# Copy customized scripts
COPY runserver-customize.sh ${ROCKETMQ_HOME}/bin/

# Expose namesrv port
EXPOSE 9876

# Override customized scripts for namesrv
# Export Java options
# Add ${JAVA_HOME}/lib/ext as java.ext.dirs
RUN mv ${ROCKETMQ_HOME}/bin/runserver-customize.sh ${ROCKETMQ_HOME}/bin/runserver.sh \
 && chmod a+x ${ROCKETMQ_HOME}/bin/runserver.sh \
 && chmod a+x ${ROCKETMQ_HOME}/bin/mqnamesrv \
 && export JAVA_OPT=" -Duser.home=/opt" \
 && sed -i 's/${JAVA_HOME}\/jre\/lib\/ext/${JAVA_HOME}\/jre\/lib\/ext:${JAVA_HOME}\/lib\/ext/' ${ROCKETMQ_HOME}/bin/tools.sh

WORKDIR ${ROCKETMQ_HOME}/bin

CMD ["/bin/bash", "mqnamesrv"]
