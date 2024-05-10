#!/bin/bash
# Make scanner
mkdir -p /app/scannerbuild/
cp -a /vagrant/code/scanner/. /app/scannerbuild/
cd /app/scannerbuild/
sleep 5
make
sleep 5
mkdir -p /app/scanner/
cp /app/scannerbuild/bin/scanner /app/scanner/scanner
rm -f /app/scanner/.env
echo "DSN = "$(cat /app/infra/docker-compose.yml | grep -o -P '(?<=DATABASE_URI: ).*(?)') >> /app/scanner/.env
sed -i 's/mariadb/127.0.0.1/g' /app/scanner/.env
echo "LOCAL_TEMP_DIR = /app/scanner/local_temp/" >> /app/scanner/.env
echo "REMOTE_WORK_DIR = /app/data/work_dir/" >> /app/scanner/.env
echo "SCANNER_IMAGE = zaproxy/zap-stable" >> /app/scanner/.env
echo "REPORTER_BIN_PATH = /app/bin/reporter" >> /app/scanner/.env
echo "LOG_LEVEL = info" >> /app/scanner/.env
echo '#!/bin/bash' > /app/bin/scanner && echo 'cd /app/scanner/ && ./scanner $@' >> /app/bin/scanner
chmod +x /app/bin/scanner

