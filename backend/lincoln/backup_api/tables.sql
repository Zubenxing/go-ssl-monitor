CREATE TABLE BackupLogs (
    Id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY, -- 自增的唯一标识符
    Ip VARCHAR(15) NOT NULL,                    -- IP 地址字段，假设 IPv4 长度最多为 15
    ServerName VARCHAR(100) NOT NULL,           -- 服务器名称
    StartTime VARCHAR(255) NOT NULL DEFAULT '',  -- 开始时间，改为 VARCHAR 类型，默认值为空字符串
    EndTime VARCHAR(255) NOT NULL DEFAULT '',    -- 结束时间，改为 VARCHAR 类型，默认值为空字符串
    BackupStatus TINYINT NOT NULL DEFAULT 0,    -- 备份状态，默认值为 0
    AlertStatus TINYINT NOT NULL DEFAULT 0 COMMENT '0:正常,1:告警已触发,2:告警未触发'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 向 BackupLogs 表中添加 ScriptVersion 字段
ALTER TABLE BackupLogs ADD COLUMN ScriptVersion VARCHAR(255) NOT NULL DEFAULT '0';