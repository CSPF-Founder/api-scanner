#!/bin/bash
#Setup Panel
mkdir -p /app/infra/
cp -a /vagrant/code/panel/. /app/infra/
#Random DB Password
sed -i "s/\[ROOT_PASS_TO_REPLACE\]/$(sudo head /dev/urandom | tr -dc 'A-Za-z0-9' | head -c 20)/g" /app/infra/docker-compose.yml
sed -i "s/\[PASSWORD_TO_REPLACE\]/$(sudo head /dev/urandom | tr -dc 'A-Za-z0-9' | head -c 20)/g" /app/infra/docker-compose.yml
cd /app/infra/
make setup
chmod 644 /app/panel/certs/panel.crt
chmod 644 /app/panel/certs/panel.key
make up
