#!/bin/bash
# Make reporter
mkdir -p /app/reporter/src
mkdir -p /app/reporter/config
rm -f /app/reporter/config/app.conf
echo "[MAIN]" >> /app/reporter/config/app.conf
echo "local_temp_dir = /app/scanner/local_temp/" >> /app/reporter/config/app.conf
echo "remote_work_dir = /app/data/work_dir/" >> /app/reporter/config/app.conf
echo "[MAIN_DATABASE]" >> /app/reporter/config/app.conf
echo "host= 127.0.0.1" >> /app/reporter/config/app.conf
echo "user= api_scanner" >> /app/reporter/config/app.conf
echo "password= "$(cat /app/infra/docker-compose.yml | grep -o -P '(?<=api_scanner:).*(?=@)') >> /app/reporter/config/app.conf
echo "db_name= api_db" >> /app/reporter/config/app.conf

echo '#!/bin/bash' > /app/bin/reporter && echo 'cd /app/reporter/src && source venv/bin/activate && python3 cli.py $@' >> /app/bin/reporter
chmod +x /app/bin/reporter

cp -a /vagrant/code/report-generator/. /app/reporter/src/
cd /app/reporter/src
apt install python3.11-venv -y
sleep 2
python3 -m venv venv
sleep 10
source venv/bin/activate
pip install poetry
poetry install --no-interaction --no-root
deactivate
exit
