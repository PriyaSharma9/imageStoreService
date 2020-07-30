#!/bin/bash
echo "Running as a docker"
./usr/share/imageStoreService/imageStoreService &
./monitor_imageStoreService.sh