#!/bin/bash
# Make manager
mkdir -p /app/managerbuild/
cp -a /vagrant/code/manager/. /app/managerbuild/
cd /app/managerbuild/
make build
mkdir -p /app/manager/
mkdir -p /app/manager/local_temp/
cp /app/managerbuild/bin/manager /app/manager/manager
chmod +x /app/manager/manager
chown -R vagrant:vagrant /app/manager/

rm -f /app/manager/.env
echo "DATABASE_URI = "$(cat /app/infra/docker-compose.yml | grep -o -P '(?<=DATABASE_URI: ).*(?)') >> /app/manager/.env
sed -i 's/mariadb/127.0.0.1/g' /app/manager/.env
echo "LOG_LEVEL = info" >> /app/manager/.env
echo "SCANNER_DOCKER = zaproxy/zap-stable" >> /app/manager/.env
echo "SCANNER_CMD = /app/bin/scanner" >> /app/manager/.env


#Making manager into service
sudo cp /app/managerbuild/apisec-manager.service /etc/systemd/system/apisec-manager.service
sudo systemctl daemon-reload
sudo systemctl enable apisec-manager.service
sudo systemctl start apisec-manager.service

