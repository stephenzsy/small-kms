FROM devopsfaith/krakend:2.4 as builder

COPY templates /etc/krakend/templates
COPY settings /etc/krakend/settings
COPY krakend.tmpl .

## Save temporary file to /tmp to avoid permission errors
RUN FC_ENABLE=1 \
    FC_OUT=/tmp/krakend.json \
    #FC_PARTIALS="/etc/krakend/partials" \
    FC_SETTINGS="/etc/krakend/settings" \
    FC_TEMPLATES="/etc/krakend/templates" \
    krakend check -d -t -c krakend.tmpl

# The linting needs the final krakend.json file
RUN krakend check -c /tmp/krakend.json --lint

FROM devopsfaith/krakend:2.4
COPY --from=builder --chown=krakend:nogroup /tmp/krakend.json .
