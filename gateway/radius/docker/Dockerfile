# Refer to https://github.com/FreeRADIUS/freeradius-server/blob/v3.2.x/scripts/docker/alpine/Dockerfile
FROM alpine:latest as build

RUN apk update
RUN apk add git gcc make
RUN apk add --no-cache --virtual .build-deps curl gpg gpg-agent

RUN mkdir -p /usr/local/src/repositories
WORKDIR /usr/local/src/repositories

RUN curl -fsSLO --compressed https://github.com/FreeRADIUS/freeradius-server/releases/download/release_3_2_3/freeradius-server-3.2.3.tar.gz \
    && curl -fsSLO https://github.com/FreeRADIUS/freeradius-server/releases/download/release_3_2_3/freeradius-server-3.2.3.tar.gz.sig \
    && curl -fsSL https://freeradius.org/pgp/packages@freeradius.org -o freeradius.gpg \
    && gpg --import freeradius.gpg
RUN gpg --verify freeradius-server-3.2.3.tar.gz.sig freeradius-server-3.2.3.tar.gz
RUN tar xzf freeradius-server-3.2.3.tar.gz

WORKDIR /usr/local/src/repositories/freeradius-server-3.2.3
RUN apk add libc-dev talloc-dev \
    openssl openssl-dev \
    linux-headers \
    pcre-dev libidn-dev krb5-dev samba-dev curl-dev json-c-dev

RUN ./configure --prefix=/opt
RUN make -j2
RUN make install
RUN rm /opt/lib/*.a

FROM alpine:latest
COPY --from=build /opt /opt

RUN apk update \
    && apk add talloc libressl pcre libwbclient tzdata \
    && apk add libcurl json-c \
    && ln -s /opt/etc/raddb /etc/raddb

RUN rm /etc/raddb/sites-available/inner-tunnel

COPY docker-entrypoint.sh /
COPY raddb/sites-enabled/default /etc/raddb/sites-enabled/default
COPY raddb/mods-enabled/eap /etc/raddb/mods-enabled/eap
COPY raddb/mods-enabled/rest /etc/raddb/mods-enabled/rest
RUN chmod +x /docker-entrypoint.sh

EXPOSE 1812/udp 1813/udp
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["radiusd"]