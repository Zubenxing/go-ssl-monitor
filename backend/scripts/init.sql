-- 创建数据库
CREATE DATABASE IF NOT EXISTS ssl_monitor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE ssl_monitor;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    last_login DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 更新默认管理员账户密码 (密码: admin123)
UPDATE users SET password = '$2a$10$JlfhGYs46sTd42mMF6keZOu8sE.X5jR6ooZ/d9J/HjD7zZ.4caOKC' WHERE username = 'admin';

-- 如果管理员账户不存在，则创建
INSERT INTO users (username, password, email, is_active)
SELECT 'admin', '$2a$10$JlfhGYs46sTd42mMF6keZOu8sE.X5jR6ooZ/d9J/HjD7zZ.4caOKC', 'admin@example.com', true
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin');

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

-- -- 创建示例数据
-- INSERT INTO domains (domain_name, notification_email, certificate_status, auto_renewal)
-- VALUES 
--     ('example.com', 'admin@example.com', 'VALID', true),
--     ('test.com', 'admin@test.com', 'VALID', true); 