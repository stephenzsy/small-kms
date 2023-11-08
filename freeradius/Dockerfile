# Refer to https://github.com/FreeRADIUS/freeradius-server/blob/v3.2.x/scripts/docker/alpine/Dockerfile
ARG from=alpine:3.18.4
FROM ${from} as build

RUN apk update
RUN apk add git gcc make
RUN apk add openssl ca-certificates
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

#
#  Install build dependencies
#
# essential
RUN apk add libc-dev talloc-dev
RUN apk add openssl openssl-dev
RUN apk add linux-headers
# general
RUN apk add pcre-dev libidn-dev krb5-dev samba-dev curl-dev json-c-dev
# RUN apk add openldap-dev unbound-dev
# languages
# RUN apk add ruby-dev perl-dev python2-dev
# databases
#RUN apk add hiredis-dev libmemcached-dev gdbm-dev libcouchbase-dev
# sql
#RUN apk add postgresql-dev mariadb-dev unixodbc-dev sqlite-dev

#
#  Build the server
#
RUN ./configure --prefix=/opt
RUN make -j2
RUN make install
RUN rm /opt/lib/*.a

#
#  Clean environment and run the server
#
FROM ${from}
COPY --from=build /opt /opt

#
# These are needed for the server to start
#
RUN apk update \
    && apk add talloc libressl pcre libwbclient tzdata \
    \
#
#  Libraries that are needed dependent on which modules are used
#  Some of these (especially the languages) are huge. A reasonable
#  selection has been enabled here. If you use modules needing
#  other dependencies then install any others required in your
#  local Dockerfile.
#
    && apk add libcurl json-c \
#libldap hiredis sqlite-dev \
#RUN apk add libidn krb5
#RUN apk add unbound-libs
#RUN apk add ruby-libs perl python2-dev
#RUN apk add libmemcached gdbm libcouchbase
#RUN apk add postgresql-dev mariadb-dev unixodbc-dev
    \
    && ln -s /opt/etc/raddb /etc/raddb
RUN apk add openssl ca-certificates

COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

EXPOSE 1812/udp 1813/udp
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["radiusd"]
