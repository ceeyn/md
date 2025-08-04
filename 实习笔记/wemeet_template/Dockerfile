FROM csighub.tencentyun.com/wemeet-public-cloud/wemeet-base-prod:latest

MAINTAINER author@tencent.com

ARG BUILD_IMAGE
ARG ENV

ENV SERVER_NAME=${BUILD_IMAGE}
ENV CONFIG_ENV=${ENV}
ENV STKE_CONTAINER_MONITOR true
ENV GOBACTRACE=crash

RUN \
    mkdir -p /home/oicq/${SERVER_NAME} &&\
    mkdir -p /home/package

COPY ./release/${SERVER_NAME} /home/oicq/${SERVER_NAME}/
RUN ls -la /home/oicq/${SERVER_NAME}/*

COPY ./sbin/entrypoint.sh /etc/kickStart.d/
RUN chmod +x /etc/kickStart.d/entrypoint.sh
