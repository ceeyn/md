#!/bin/bash

# ====================================================
# eg :
#     ./start.sh start
#
# ====================================================

# action: start stop reload restart status ...

action=$1

# root path of app
base_path=$(dirname $(cd `dirname $0`; pwd))
app_bin=$(basename ${base_path})
app_path="${base_path}/bin/${app_bin}"

echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] APP Path: ${base_path}"
echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] APP Name: ${app_bin}"
echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] APP Bin: ${app_path}"
echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] Action: ${action}"

cd ${base_path}
chmod 777 ${app_path}

function Start(){
    if [ ! -d "${base_path}/log" ]; then
        mkdir -p ${base_path}/log
    fi

    pid=`get_pid`
    if [ -n "$pid" ]; then
        echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] ${app_bin} is running! PID:${pid}"
        return;
    fi

    echo -en "$(date '+%Y-%m-%d %H:%M:%S') [INFO] Run ${app_bin} service:  "
    nohup ${app_path} >> ${base_path}/log/nohup.log 2>&1 &
    [ -n "`get_pid`" ] && echo "[ OK ]" || echo "[ FAIL ]"

    sleep 3;
    check=`get_pid`
    if [ -n "$check" ]; then
        echo -e "$(date '+%Y-%m-%d %H:%M:%S') [INFO] Start ${app_bin} service:  [ OK ]"
    else
        echo -e "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] Start ${app_bin} service:  [ FAIL ]"
    fi
}

function Stop(){
    pid=`get_pid`
    if [ -n "${pid}" ]; then
        kill ${pid}
        echo -e "$(date '+%Y-%m-%d %H:%M:%S') [INFO] Kill ${app_bin} service:  [ OK ]"
    else
        echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] ${app_bin} is not running!"
        return;
    fi
    for((i=1;i<=100;i++));
    do
        sleep 0.5;
        check=`get_pid`
        if [ -z "$check" ]; then
          echo -e "$(date '+%Y-%m-%d %H:%M:%S') [INFO] Stop ${app_bin} service:  [ OK ]"
          return;
        fi
    done
    echo -e "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] Stop ${app_bin} service:  [ FAIL ]"
}

function Reload(){
    pid=`get_pid`
    if [ -n "${pid}" ]; then
        kill -HUP ${pid}
        echo -e "$(date '+%Y-%m-%d %H:%M:%S') [INFO] Reload ${app_bin} service:  [ OK ]"
    else
        echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] ${app_bin} is not running!"
    fi
}

function Status(){
    pid=`get_pid`
    if [ -n "$pid" ]; then
        echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] ${app_bin} is running! PID:${pid}"
    else
        echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] ${app_bin} is not running!"
    fi
}

function get_pid(){
    pid="`ps -f -u $(whoami) | grep ${app_path} | grep -v 'grep' | sort -r | awk '{print $2;exit}'`"
    echo $pid
    return $pid
}

function Usage(){
    echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] there must be 1 param ..."
    echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] $0 <start|stop|restart|status|reload>"
}

if [ $# -lt 1 ]; then
    Usage
    exit 0
else
	# assert bin file
	if [ ! -e "${app_path}" ]; then
	    echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] ${app_path} file not exists. ${action} failed"
	    exit 1
	fi

	if [ ! -x "${app_path}" ]; then
	    echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] ${app_path} without execute permission. ${action} failed"
	    exit 1
	fi

	case $1 in
	    start)
	        Start
	        ;;
	    stop)
	        Stop
	        ;;
	    status)
	        Status
	        ;;
	    restart)
	        Stop
	        Start
	        ;;
	    reload)
	        Reload
	        ;;
	    *)
	        Usage
	        ;;
	esac
fi