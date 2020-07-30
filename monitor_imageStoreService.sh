#!/bin/bash
while true
do
  sleep 60
  now=$(date)
  PROCESS_NUM=$(ps -ef | grep "\.*/imageStoreService" | grep -v "grep" | grep -v "vi" | grep -v build | grep -v "tail"| grep -v "bash" | wc -l)

  if [ $PROCESS_NUM -eq 0 ];
  then
     echo "ImageStoreService is not running, starting now $now" >> /monitor_imageStoreService.log 2>&1
     nohup ./usr/share/imageStoreService/imageStoreService &
  fi
done