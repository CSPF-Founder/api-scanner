#!/bin/bash
#Install Golang
cd ~
wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
sudo echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile
