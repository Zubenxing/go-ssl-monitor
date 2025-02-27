#!/bin/bash
#2025-02-25
#lincolnzhang@techtrex.com
#version: 1.0.0.1

script_version="1.0.0.1"

apiServerAdd="http://20.83.180.26:8080"

rsync_name=$(/usr/bin/hostname)
rsync_user=$(/usr/bin/hostname | awk -F '-' '{print$2}')
logs_dir="/back/rsync/logs"
backupServer="backup.kiosoft.com"
#log month
month=$(/usr/bin/date '+%Y%m')
#cpu limit percentage by cpu cores
percentage=$(($(nproc) * 10))

public_ip=$(curl -s https://ipinfo.io/ip) #获取公网IP

#log function
function log() {
    log="$(date '+%Y-%m-%d %H:%M:%S') $@"
    echo "$log" >>/var/log/rsync-"$month".log
}
echo -e >>/var/log/rsync-"$month".log
log "cpulimit $percentage"

#=================================================================================
# 函数：压缩日志并移动到目标目录
# 参数：
#   $1：源目录路径
#   $2：目标目录路径
#   例:
#       compress_and_move_logs /var/log/nginx /var/backup/logs
#=================================================================================

function compress_and_move_logs() {
    #判断传入的值不为空
    if [ -z "$1" ]; then
        log '$1 参数为空，退出函数'
        return 1 # 返回1表示函数退出
    fi
    if [ -z "$2" ]; then
        log '$2 参数为空，退出函数'
        return 1 # 返回1表示函数退出
    fi

    local src_dir="$1"                   # 源目录 / Source directory
    local dest_dir="/back/rsync/logs/$2" # 目标目录 / Destination directory
    #判断目录是否正在
    # 检查源目录和目标目录是否存在
    # Check if source and destination directories exist
    if [ ! -d "$src_dir" ]; then
        log "源目录不存在: $src_dir / Source directory does not exist"
        return 1
    fi

    if [ ! -d "$dest_dir" ]; then
        log "目标目录不存在创建目录: $dest_dir / Destination directory does not exist"
        mkdir -p "$dest_dir"
    fi
    # 获取前一天的日期 (格式：YYYY-MM-DD)
    # Get the previous day's date (format: YYYY-MM-DD)
    yesterday=$(date -d "yesterday" +"%Y-%m-%d")
    yesterdayFormant=$(date -d "yesterday" +"%Y%m%d")

    # 查找源目录中文件名包含前一天日期的 .log 文件
    # Find .log files in the source directory that contain the previous day's date

    find "$src_dir" -type f \( -name "*${yesterday}.log" -o -name "log-${yesterday}.php" -o -name "*${yesterday}.log.locking" -o -name "*.log.${yesterday}" -o -name "*${yesterdayFormant}.log" \) | while read -r log_file; do
        # 压缩文件 / Compress the file
        if [ -f "$log_file" ]; then
            gzip --rsyncable "$log_file"
            # 检查压缩是否成功 / Check if compression was successful
            if [ -f "${log_file}.gz" ]; then
                # 移动压缩后的文件到目标目录 / Move compressed file to destination
                mv "${log_file}.gz" "$dest_dir"
                log "已压缩并移动文件: ${log_file}.gz / Compressed and moved: ${log_file}.gz"
            else
                log "压缩文件失败: $log_file / Compression failed"
            fi

        else
            continue
        fi
    done

    #搜索压缩完的日志移动到dest_dir目录

}


#=================================================================================
#判断上一次备份是否成功，备份成功则删除昨天的日志不成功则不删除
#
#================================================================
response=$(curl -s "$apiServerAdd/getStatus/$public_ip")
backup_status=$(echo "$response" | grep -oP '"backup_status":\K\d+')
echo $response
if [ -n "$backup_status" ]; then
    log "获取到的 status: $backup_status"
    if [ $backup_status -eq 0 ]; then
        find /back/rsync -type f -name "*.gz" -exec rm -f {} \;
    fi
else
    log "未能获取 backup_status,不删除"
fi

#============================================================================
#记录开始时间
#============================================================================
curl -X POST $apiServerAdd/backupLogs \
    -H "Content-Type: application/json" \
    -d "{
    \"ip\": \"$public_ip\",
    \"server_name\": \"$rsync_name\",
    \"start_time\": \"$(date '+%Y-%m-%d %H:%M:%S')\",
    \"backup_status\": 1,
    \"alert_status\": 1,
    \"script_version\": \"$script_version\"
}"
if [ $? -ne 0 ]; then
    log "Failed to post backup start log"
    exit 1
fi
response=$(curl -s "$apiServerAdd/getId/$public_ip")
id=$(echo "$response" | grep -oP '"id":\K\d+')
if [ -n "$id" ]; then
    log "获取到的 ID: $id"
else
    log "未能获取 ID"
fi

#============================================================================
# Compress Azure files and perform backups.
#============================================================================
function backup_azure_file() {
    # Check if the parameters are empty
    if [ -z "$1" ] || [ -z "$2" ]; then
        echo "Error: Both source and destination directories are required."
        return 1
    fi

    # Check if the backup directory under the destination directory exists. If not, create it.
    local backup_dir="$2/$(basename "$1")"
    if [ ! -d "$backup_dir" ]; then
        mkdir -p "$backup_dir"
        if [ $? -eq 0 ]; then
            echo "create path $backup_dir"
        else
            echo "failed to create path $backup_dir"
            return 1
        fi
    fi

    # Generate the backup file name using a date format without spaces
    local backup_file="$backup_dir/$(basename "$1")_$(date '+%Y%m%d%H%M%S').tar.gz"
    # Perform the backup operation
    tar zcf "$backup_file" "$1"
    if [ $? -eq 0 ]; then
        echo "Backup successful: $backup_file"
    else
        echo "Backup failed: $backup_file"
        return 1
    fi
}
#============================================================================
#开始查找日志并压缩日志
#============================================================================
#laundry portal
compress_and_move_logs "/KiosoftApplications/WebApps/kiosk_laundry_portal/application/logs" "laundry_portal"
#value code
compress_and_move_logs "/KiosoftApplications/WebApps/kiosk_value_code/application/logs" "value_code"
#web_lcms
compress_and_move_logs "/KiosoftApplications/WebApps/kiosk_web_lcms/application/logs" "web_lcms"
#web_rss
compress_and_move_logs "/KiosoftApplications/WebApps/kiosk_web_rss/application/logs" "web_rss"
#ccm_bridge
compress_and_move_logs "/KiosoftApplications/ServerApps/TTI_ccm_bridge/logs" "ccm_bridge"
#proxy
compress_and_move_logs "/KiosoftApplications/ServerApps/TTI_Proxy/logs" "proxy"
#reportserver
compress_and_move_logs "/KiosoftApplications/ServerApps/TTI_ReportServer/logs" "reportServer"
#sms_python
compress_and_move_logs "/KiosoftApplications/ServerApps/TTI_sms_python/app/logs" "sms_python"
#tcp
compress_and_move_logs "/KiosoftApplications/ServerApps/TTI_tcp/logs" "tcp"
#universal_ddd
compress_and_move_logs "/KiosoftApplications/ServerApps/TTI_universal_ddd/logs" "universal_ddd"
#rabbitmq
compress_and_move_logs "/KiosoftApplications/ServerApps/rabbitmq_api/logs" "rabbitmq"
#token server
compress_and_move_logs "/tmp/TokenServerLog" "tokenServer"
#policy server
compress_and_move_logs "/tmp/policyServerLog" "policyServer"
#universal_ddd
compress_and_move_logs "/KiosoftApplications/WebApps/universal_ddd/storage/logs" "universal_ddd"
#rabbitmq_api
compress_and_move_logs "/KiosoftApplications/ServerApps/rabbitmq_api/logs" "rabbitmq_api"
#coffee
compress_and_move_logs "/KiosoftApplications/WebApps/kiosk_coffee_back_end/storage/logs" "coffee"
#kps
compress_and_move_logs "/KiosoftApplications/WebApps/kps_backend/storage/logs" "kps"

#Example bakcup azure file
#backup_azure_file "/mnt/sa25omdryrun01/fs01-sa25omdryrun01" "/back/rsyn"

#============================================================================
#随机数函数准备开始同步
#============================================================================
function randsecond() {
    min=$1
    max=$(($2 - $min + 1))
    num=$(date +%s%N)
    echo $(($num % $max + $min))
}

second=$(randsecond 1 1800)
log "Sleep for $second seconds"
sleep "$second"

#write time to log
day=$(/usr/bin/date '+%Y-%m-%d %H:%M:%S')
echo -e >>/var/log/rsync-"$month".log
echo "$day" >>/var/log/rsync-"$month".log

#rsync to remote backup server
log "Start sending"
#rsync to remote backup server
cpulimit -l $percentage /usr/bin/rsync -avz --password-file=/etc/rsync.password /back/rsync/* "$rsync_user"@$backupServer::"$rsync_name" >>/var/log/rsync-"$month".log
if [ $? -eq 0 ]; then
    log "rsync 成功"

    curl -X PUT "$apiServerAdd/backupLogs/$id" \
        -H "Content-Type: application/json" \
        -d "{
    \"end_time\": \"$(date '+%Y-%m-%d %H:%M:%S')\",
    \"backup_status\": 0,
    \"alert_status\": 0
}"

else
    echo "rsync 失败"

    curl -X PUT "$apiServerAdd/backupLogs/$id" \
        -H "Content-Type: application/json" \
        -d "{
    \"end_time\": \"$(date '+%Y-%m-%d %H:%M:%S')\",
    \"backup_status\": 1,
    \"alert_status\": 1
}"

fi

log "Backup complete"
exit 0


