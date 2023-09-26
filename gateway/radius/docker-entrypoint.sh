#!/bin/sh
set -e

cp /run/secrets/radius-tls-key.pem /etc/raddb/certs/radius-tls-key.pem
cp /run/secrets/radius-tls-server.pem /etc/raddb/certs/radius-tls-server.pem
cp /run/secrets/radius-tls-ca.pem /etc/raddb/certs/radius-tls-ca.pem

chmod 0400 /etc/raddb/certs/radius-tls-key.pem /etc/raddb/certs/radius-tls-server.pem /etc/raddb/certs/radius-tls-ca.pem

PATH=/opt/sbin:/opt/bin:$PATH
export PATH

# this if will check if the first argument is a flag
# but only works if all arguments require a hyphenated flag
# -v; -SL; -f arg; etc will work, but not arg1 arg2
if [ "$#" -eq 0 ] || [ "${1#-}" != "$1" ]; then
    set -- radiusd "$@"
fi

# check for the expected command
if [ "$1" = 'radiusd' ]; then
    shift
    exec radiusd -f "$@"
fi

# debian people are likely to call "freeradius" as well, so allow that
if [ "$1" = 'freeradius' ]; then
    shift
    exec radiusd -f "$@"
fi

# else default to run whatever the user wanted like "bash" or "sh"
exec "$@"

