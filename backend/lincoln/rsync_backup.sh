#!/bin/bash
#2024-11-13
#lincolnzhang@techtrex.com

#globle variable
webapp="/KiosoftApplications/WebApps"
serverapp="/KiosoftApplications/ServerApps"
logs_dir="/back/rsync/logs"
code_dir="/back/rsync/code"
apiServerAdd="http://20.83.180.26:8080"



rsync_name=$(/usr/bin/hostname)
rsync_user=$(/usr/bin/hostname | awk -F '-' '{print$2}')

backupServer="backup.kiosoft.com"
#log month
month=$(/usr/bin/date '+%Y%m')
#cpu limit percentage by cpu cores
percentage=$(($(nproc) * 10))

public_ip=$(curl -s https://ipinfo.io/ip) #获取公网IP

#log function
function log()
{
    log="$(date '+%Y-%m-%d %H:%M:%S') $@"
    echo $log >> /var/log/rsync-$month.log
}
echo -e >> /var/log/rsync-$month.log
log "cpulimit $percentage"




log "Starting pack ..."
#create dir
if [ ! -d "$logs_dir/universal_ddd" ] ;then
    mkdir -p /back/rsync/{code,database,logs}
    mkdir -p $logs_dir/{universal_ddd,rabbitmq}
fi

#1.code backup
log "start pack code"
date_minute=`date +%Y%m%d%H%M`
rm -rf $code_dir/*
#====================================================================================================================
cpulimit -l $percentage /usr/bin/rsync -avzrtopg $webapp/universal_ddd $code_dir --exclude '*.log' --exclude '*.gz'
cpulimit -l $percentage /usr/bin/rsync -avzrtopg $webapp/universal_frontend_build $code_dir
cpulimit -l $percentage /usr/bin/rsync -avzrtopg $webapp/powerbi_laundry $code_dir
cpulimit -l $percentage /usr/bin/rsync -avzrtopg $webapp/powerbi_coffee $code_dir
cd $code_dir
cpulimit -l $percentage zip -r web_code_$date_minute.zip *
find $code_dir -maxdepth 1 -mindepth 1 -type d | xargs -i rm -rf {}
log "pack code complete"

#2.log
#清空目录，不然会有重复压缩问题
log "start pack log"

#=================================================================================
# 函数：压缩日志并移动到目标目录
# Function: Compress log files and move them to the destination directory
universal_DDD_compress_and_move_logs() {
    local src_dir="$1"   # 源目录 / Source directory
    local dest_dir="$2"  # 目标目录 / Destination directory

    # 检查源目录和目标目录是否存在
    # Check if source and destination directories exist
    if [ ! -d "$src_dir" ]; then
        echo "源目录不存在: $src_dir / Source directory does not exist"
        return 1
    fi

    if [ ! -d "$dest_dir" ]; then
        echo "目标目录不存在: $dest_dir / Destination directory does not exist"
        return 1
    fi

    # 获取前一天的日期 (格式：YYYY-MM-DD)
    # Get the previous day's date (format: YYYY-MM-DD)
    yesterday=$(date -d "yesterday" +"%Y-%m-%d")

    # 查找源目录中文件名包含前一天日期的 .log 文件
    # Find .log files in the source directory that contain the previous day's date
    find "$src_dir" -type f -name "*${yesterday}*.log" | while read -r log_file; do
        # 压缩文件 / Compress the file
        gzip --rsyncable "$log_file"
        
        # 检查压缩是否成功 / Check if compression was successful
        if [ -f "${log_file}.gz" ]; then
            # 移动压缩后的文件到目标目录 / Move compressed file to destination
            mv "${log_file}.gz" "$dest_dir"
            log "已压缩并移动文件: $(basename "${log_file}.gz") / Compressed and moved: $(basename "${log_file}.gz")"
        else
            log "压缩文件失败: $log_file / Compression failed"
        fi
    done
}


# 函数：压缩日志并移动到目标目录
# Function: Compress log files and move them to the destination directory
rabbitmq_compress_and_move_logs() {
    local src_dir="$1"   # 源目录 / Source directory
    local dest_dir="$2"  # 目标目录 / Destination directory

    # 检查源目录和目标目录是否存在
    # Check if source and destination directories exist
    if [ ! -d "$src_dir" ]; then
        echo "源目录不存在: $src_dir / Source directory does not exist"
        return 1
    fi

    if [ ! -d "$dest_dir" ]; then
        echo "目标目录不存在: $dest_dir / Destination directory does not exist"
        return 1
    fi

    # 获取前一天的日期 (格式：YYYY-MM-DD)
    # Get the previous day's date (format: YYYY-MM-DD)
    yesterday=$(date -d "yesterday" +"%Y-%m-%d")

    # 查找源目录中文件名包含前一天日期的 .log 文件
    # Find .log files in the source directory that contain the previous day's date
      find "$src_dir" -type f -name "rabbitmq_*.log.${yesterday}" | while read -r log_file; do
        # 压缩文件 / Compress the file
        gzip --rsyncable "$log_file"
        
        # 检查压缩是否成功 / Check if compression was successful
        if [ -f "${log_file}.gz" ]; then
            # 移动压缩后的文件到目标目录 / Move compressed file to destination
            mv "${log_file}.gz" "$dest_dir"
            log "已压缩并移动文件: $(basename "${log_file}.gz") / Compressed and moved: $(basename "${log_file}.gz")"
        else
            log "压缩文件失败: $log_file / Compression failed"
        fi
    done
}

#=================================================================================
#判断上一次备份是否成功，成功删除昨天的日志不成功不删除
#
#
#rm -rf $logs_dir/universal_ddd/*
#rm -rf $logs_dir/rabbitmq/*
#================================================================
response=$(curl -s "$apiServerAdd/getStatus/$public_ip")
backup_status=$(echo "$response" | grep -oP '"backup_status":\K\d+')
echo $response
if [ -n "$backup_status" ]; then
    log "获取到的 status: $backup_status"
    if [ $backup_status -eq 0 ]; then
        find $logs_dir/universal_ddd/ -type f -name "*.gz"  -exec rm -f {} \;
        find $logs_dir/rabbitmq/ -type f -name "*.gz" -exec rm -f {} \;
    fi
else
    log "未能获取 backup_status,不删除"
fi

#============================================================================
#记录执行时间
#============================================================================
curl -X POST $apiServerAdd/backupLogs \
-H "Content-Type: application/json" \
-d "{
    \"ip\": \"$public_ip\",
    \"server_name\": \"$rsync_name\",
    \"start_time\": \"$(date '+%Y-%m-%d %H:%M:%S')\",
    \"backup_status\": 1,
    \"alert_status\": 1
}"

response=$(curl -s "$apiServerAdd/getId/$public_ip")
id=$(echo "$response" | grep -oP '"id":\K\d+')
if [ -n "$id" ]; then
    log "获取到的 ID: $id"
else
    log "未能获取 ID"
fi

#============================================================================
#find $webapp/universal_ddd/storage/logs/ -type f -name "*.log" -mtime +3 -exec mv {} $logs_dir/universal_ddd/ \;
#cpulimit -l $percentage /usr/bin/rsync -avzrtopg `find $webapp/universal_ddd/storage/logs/ -type f -name "*.log"` $logs_dir/universal_ddd/
#cpulimit -l $percentage /usr/bin/rsync -avzrtopg $(find $serverapp/rabbitmq_api/logs/ -type f -name "*.log.*") $logs_dir/rabbitmq/
#============================================================================


universal_DDD_compress_and_move_logs "$webapp/universal_ddd/storage/logs" "$logs_dir/universal_ddd"
rabbitmq_compress_and_move_logs "$serverapp/rabbitmq_api/logs" "$logs_dir/rabbitmq"


#gzip logs========================================================================
find $logs_dir -type f ! -name "*.gz" | xargs -i cpulimit -l $percentage gzip --rsyncable {}
log "pack log complete"

#Wait for random seconds to avoid server congestion
function randsecond(){
  min=$1
  max=$(($2-$min+1))
  num=$(date +%s%N)
  echo $(($num%$max+$min))
}

second=$(randsecond 1 3600)
log "Sleep for $second seconds"
sleep $second

#write time to log
month=`/usr/bin/date '+%Y%m'`
day=`/usr/bin/date '+%Y-%m-%d %H:%M:%S'`
echo -e >> /var/log/rsync-$month.log
echo $day >> /var/log/rsync-$month.log

#rsync to remote backup server
log "Start sending"
#rsync to remote backup server
cpulimit -l $percentage /usr/bin/rsync -avz --password-file=/etc/rsync.password /back/rsync/* $rsync_user@$backupServer::$rsync_name >> /var/log/rsync-$month.log
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



