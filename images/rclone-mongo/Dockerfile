FROM --platform=$TARGETPLATFORM mongo:4.2.24 as all

USER root

COPY rootfs/ /

# To fix the error: The following signatures were invalid: EXPKEYSIG 4B7C549A058F8B6B MongoDB 4.2 Release Signing Key <packaging@mongodb.com>
RUN apt-get update --allow-insecure-repositories

RUN apt-get install -y unzip curl bash

RUN curl https://rclone.org/install.sh | bash

RUN cp /usr/bin/rclone /usr/local/bin/

RUN cp /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime

RUN chmod +x /usr/local/bin/rclone

RUN chmod +x /usr/local/bin/docker-entrypoint.sh