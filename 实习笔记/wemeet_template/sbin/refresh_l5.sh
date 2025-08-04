#!/bin/sh

# ====================================================
#
#  L5预热脚本
#
# ====================================================

# 需要预热的CL5文件
conf_file_array=("conf/trpc_go.yaml")


base_path=$(dirname $(cd `dirname $0`; pwd))

for element in ${conf_file_array[*]};
do
    #target: cl5://2054401:65536
    #target = "cl5://2054401:65536"
    while read line
    do
        if [[ $line =~ "cl5" ]]; then
            TMP1=${line##*//}
            TMP2=${TMP1%\"*}
            TMP3=${TMP2//:/ }
            TMP4=${TMP3%#*}
            echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] ${line} refresh cl5 [$TMP4] ..."
            /usr/local/services/l5_protocol_32os-1.0/bin/L5GetRoute1 $TMP4 5
        fi
    done <${base_path}/${element}
done
