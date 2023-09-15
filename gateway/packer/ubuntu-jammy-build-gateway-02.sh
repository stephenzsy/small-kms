#!/bin/bash -e

az login --identity
az acr login --name $ACR_NAME
docker pull $ACR_IMAGE_RADIUS

# sync
/usr/sbin/waagent -force -deprovision+user && export HISTSIZE=0 && sync
