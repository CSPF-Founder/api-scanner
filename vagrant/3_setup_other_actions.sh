#!/bin/bash

#Install other depedencies
sudo apt-get install make -y 


#Make app dirs
mkdir -p /app/bin/
mkdir -p /app/panel/
chown -R vagrant:vagrant /app/
#Other Depedencies
sudo apt install net-tools sudo wget nano telnet python3-pip python3-poetry python3-venv  ntpdate python3.11-venv -y

#Set time
sudo ntpdate pool.ntp.org


# Other configs
sudo sed -i 's/#SystemMaxUse=/SystemMaxUse=100M/g' /etc/systemd/journald.conf

