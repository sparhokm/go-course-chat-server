FROM alpine:3.19

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /usr/local/bin/goose
RUN chmod +x /usr/local/bin/goose

COPY ./common/wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

COPY ./development/migrator/migrator.sh /usr/local/bin/migration
RUN chmod +x /usr/local/bin/migration

WORKDIR /root
