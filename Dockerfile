FROM registry.cn-shanghai.aliyuncs.com/zhangju/alpine:3.7

ENV INSTALL_STATUS=success \
    ENV_CONFIG=open \
    DB_ENGINE=mysql \
    DB_HOST=127.0.0.1 \
    DB_USER=root \
    DB_PASSWORD=root \
    DB_PORT=3600 \
    DB_DATABASE=cron

# 设置源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update

RUN apk add --no-cache ca-certificates tzdata g++ gcc \
    && addgroup -S app \
    && adduser -S -g app app

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

COPY ./bin /app/

RUN mkdir /app/log && chown app:app -R /app && chmod -R 777 /app

EXPOSE 5920

USER app

VOLUME ["/app/log"]

ENTRYPOINT ["/app/gocron","web"]