#!/bin/bash

# ====================================================
#
#  进程拉起脚本，监控到进程小时，自动拉起服务
#
# ====================================================

base_path=$(dirname $(cd `dirname $0`; pwd))
app_bin=$(basename ${base_path})
app_path="${base_path}/bin/${app_bin}"

pid="`ps -f -u $(whoami) | grep ${app_path} | grep -v 'grep' | sort -r | awk '{print $2;exit}'`"
if [ -n "$pid" ]; then
    echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] ${app_bin} is running! PID: ${pid}"
else
    echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] ${app_bin} is not running ..."
    su root -c "${base_path}/sbin/start.sh start >> ${base_path}/log/start.log"
fi
