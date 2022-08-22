#!/bin/bash

#每分钟检测gofound运行
#*/1 * * * * /data/gofound/gofound.sh > /dev/null 2>&1

count=`ps -fe |grep "gofound"|grep "config.yaml" -c`

echo $count
if [ $count -lt 1 ]; then
	echo "restart"
	echo $(date +%Y-%m-%d_%H:%M:%S) >/data/gofound/restart.log 
	/etc/init.d/gofound.d restart
else
	echo "is running"
fi