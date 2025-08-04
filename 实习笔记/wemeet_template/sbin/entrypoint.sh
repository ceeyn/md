#!/bin/bash

# ====================================================
#
#  入口脚本 (2021.05.12)
#
# ====================================================

base_path="/home/oicq/${SERVER_NAME}"
if [ ! -d "${base_path}" ]; then
    echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] dir: ${base_path} not found ... " >> ${base_path}/log/entrypoint.log
    exit 1;
fi

echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] entrypoint run start ... " >> ${base_path}/log/entrypoint.log

#建立日志软连接
mkdir -p ${base_path}/logs/
rm -rf ${base_path}/logs/
mkdir -p /app/logs/${SERVER_NAME}/
ln -nsf /app/logs/${SERVER_NAME}/ ${base_path}/logs

cd ${base_path};

# 预执行环境变量注入
su root -c "${base_path}/sbin/injection_var.sh >> ${base_path}/log/injection_var.log" || exit 1;

# 预执行刷新l5脚本
su root -c "${base_path}/sbin/refresh_l5.sh >> ${base_path}/log/refresh_l5.log"

# 预执行刷新cmlb脚本
su root -c "${base_path}/sbin/refresh_cmlb.sh >> ${base_path}/log/refresh_cmlb.log"


# l5脚本 加入 crontab
echo "*/3 * * * * ${base_path}/sbin/refresh_l5.sh >> ${base_path}/log/refresh_l5.log 2>&1" >> /var/spool/cron/root

# cmlb脚本 加入 crontab
echo "*/3 * * * * ${base_path}/sbin/refresh_cmlb.sh >> ${base_path}/log/refresh_cmlb.log 2>&1" >> /var/spool/cron/root

# 拉起脚本 加入 crontab （5s 一次）
for((i=0;i<60;i+=5));
do
  echo "*/1 * * * * sleep ${i}s && ${base_path}/sbin/monitor_process.sh >> ${base_path}/log/monitor_process.log 2>&1" >> /var/spool/cron/root
done

# 日志清理脚本 加入 crontab
echo "*/1 * * * * ${base_path}/sbin/clean_log.sh >> ${base_path}/log/clean_log.log 2>&1" >> /var/spool/cron/root

# 启动程序
su root -c "${base_path}/sbin/start.sh start >> ${base_path}/log/start.log"

echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] entrypoint run finish ... " >> ${base_path}/log/entrypoint.log
