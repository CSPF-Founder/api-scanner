#!/bin/bash
#Download catz
cd ~
wget https://github.com/Endava/cats/releases/download/cats-8.6.0/cats_linux_amd64_8.6.0.tar.gz
sudo tar -C /usr/local -xzf cats_linux_amd64_8.6.0.tar.gz
sudo mv /usr/local/cats /app/bin/qats
chmod +x /app/bin/qats
