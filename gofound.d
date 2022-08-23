#!/bin/bash
# chkconfig: 2345 90 10
# Description: Startup script for gofound on Debian. Place in /etc/init.d and
# run 'update-rc.d -f gofound defaults', or use the appropriate command on your
# distro. For CentOS/Redhat run: 'chkconfig --add gofound'

### BEGIN INIT INFO
#  
# Provides:         gofound.d
# Required-Start:   $local_fs  $remote_fs  
# Required-Stop:    $local_fs  $remote_fs  
# Default-Start:    2 3 4 5
# Default-Stop:     0 1 6
# Short-Description: starts gofound
# Description:       This file should be used to gofound scripts to be placed in /etc/init.d.  
#  
### END INIT INFO 


## 2345是默认启动级别，级别有0-6共7个级别。 90是启动优先级，10是停止优先级，优先级范围是0－100，数字越大，优先级越低。

## Fill in name of program here.  
PROG="gofound"
PROG_PATH="/usr/local/bin" ## Not need, but sometimes helpful (if $PROG resides in /opt for example). 
PROG_ARGS="--config=/gofound_path/config.yaml"
PID_PATH="/var/run/"
  
start() {  
    if [ -e "$PID_PATH/$PROG.pid" ]; then  
        ## Program is running, exit with error.  
        echo "Error! $PROG_PATH/$PROG is currently running!" 1>&2
        exit 1  
    else  
        ## Change from /dev/null to something like /var/log/$PROG if you want to save output.  
        $PROG_PATH/$PROG $PROG_ARGS 2>&1 >>/var/log/$PROG &
        #pid=`ps ax | grep -i '/usr/bin/frps' | grep -v 'grep' |  sed 's/^\([0-9]\{1,\}\).*/\1/g' | head -n 1`  
        pid=`ps -ef | grep $PROG_PATH/$PROG | grep -v grep | awk '{print $2}'`
        #echo $PROG_PATH/$PROG $PROG_ARGS 
        echo "$PROG_PATH/$PROG($pid) started"  
        echo $pid > "$PID_PATH/$PROG.pid"  
    fi  
}  
  
stop() {  
    echo "begin stop"  
    if [ -e "$PID_PATH/$PROG.pid" ]; then  
        ## Program is running, so stop it  
        #pid=`ps ax | grep -i '/usr/bin/frps' | grep -v 'grep' | sed 's/^\([0-9]\{1,\}\).*/\1/g' | head -n 1`  
        pid=`ps -ef | grep $PROG_PATH/$PROG | grep -v grep | awk '{print $2}'`
        kill $pid  

        rm -f  "$PID_PATH/$PROG.pid"  
        echo "$PROG_PATH/$PROG($pid) stopped"
    else  
        ## Program is not running, exit with error.  
        echo "Error! $PROG_PATH/$PROG not started!" 1>&2
    fi  
}

status() {  
    if [ -e "$PID_PATH/$PROG.pid" ]; then  
        ## Program is running, so stop it  
        #pid=`ps ax | grep -i '/usr/bin/frps' | grep -v 'grep' | sed 's/^\([0-9]\{1,\}\).*/\1/g' | head -n 1`  
        pid=`ps -ef | grep $PROG_PATH/$PROG | grep -v grep | awk '{print $2}'`
 
		if [ $pid ]; then
			echo "$PROG_PATH/$PROG($pid) is running..."
		else
			echo "$PROG_PATH/$PROG dead but pid file exists" 1>&2
		fi
    else  
        ## Program is not running, exit with error.  
        echo "Error! $PROG_PATH/$PROG not started!" 1>&2
    fi  
}

  
## Check to see if we are running as root first.  
## Found at http://www.cyberciti.biz/tips/shell-root-user-check-script.html  
if [ "$(id -u)" != "0" ]; then  
    echo "This script must be run as root" 1>&2
    exit 1  
fi  
  
case "$1" in
    start)
        start
        exit 0
    ;;
    stop)
        echo '' > /var/log/$PROG
        stop
        exit 0
    ;;  
    reload|restart|force-reload)
        stop
        start
        exit 0
    ;;
    status)
		status
        exit 0
	;;
    *)
        echo "Usage: $0 {start|stop|restart|status}" 1>&2
        exit 1
    ;;
esac
