#!/bin/bash

# Chmod shared libs

mkdir -p /app/data/
mkdir -p /app/scanner/local_temp/
sudo chown -R vagrant:vagrant /app/data/
sudo chown -R vagrant:vagrant /app/scanner/
sudo chown -R vagrant:vagrant /app/reporter/

chmod -R 777 /app/data/
chmod -R 777 /app/scanner/
chmod -R 777 /app/reporter/

# Cleanup
sudo rm -rf /app/scannerbuild/
sudo rm -rf /app/managerbuild/