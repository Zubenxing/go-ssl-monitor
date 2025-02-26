-- 创建数据库
CREATE DATABASE IF NOT EXISTS ssl_monitor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE ssl_monitor;

-- 创建domains表
CREATE TABLE IF NOT EXISTS domains (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    domain_name VARCHAR(255) NOT NULL UNIQUE,
    notification_email VARCHAR(255),
    certificate_status VARCHAR(50),
    certificate_issuer VARCHAR(255),
    certificate_expiry_date DATETIME,
    last_checked DATETIME,
    auto_renewal BOOLEAN DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_domain_name (domain_name),
    INDEX idx_certificate_status (certificate_status),
    INDEX idx_certificate_expiry_date (certificate_expiry_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建示例数据
INSERT INTO domains (domain_name, notification_email, certificate_status, auto_renewal)
VALUES 
    ('example.com', 'admin@example.com', 'VALID', true),
    ('test.com', 'admin@test.com', 'VALID', true);

-- 创建用户表（预留）
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 