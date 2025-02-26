# SSL证书监控系统

一个基于Go和Vue3的SSL证书监控系统，用于监控多个域名的SSL证书状态和自动续期。

## 功能特点

- 多域名SSL证书监控
- 证书到期提醒
- 自动续期功能（基于Let's Encrypt）
- 友好的Web界面
- 定时检查证书状态
- 邮件通知功能

## 技术栈

### 后端 (Go)
- Go 1.21+
- Gin Web Framework
- GORM
- MySQL 8.0
- JWT认证
- 定时任务调度

### 前端
- Vue 3
- Element Plus
- Axios
- Vue Router
- Pinia
- TypeScript
- Vite

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0
- Redis (可选，用于缓存)

### 数据库配置
1. 创建MySQL数据库：
```sql
CREATE DATABASE ssl_monitor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 初始化数据库表：
```bash
cd backend
mysql -h your_host -u your_user -p your_database < scripts/init.sql
```

### 后端启动
1. 配置数据库连接（backend/configs/config.yaml）：
```yaml
mysql:
  host: "your_host"
  port: 3306
  user: "your_user"
  password: "your_password"
  database: "ssl_monitor"
```

2. 启动后端服务：
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 前端启动
```bash
cd frontend/vue-project
npm install
npm run dev
```

## API文档

### 域名管理API
- GET /api/domains - 获取所有域名
- POST /api/domains - 添加新域名
- PUT /api/domains/:id - 更新域名信息
- DELETE /api/domains/:id - 删除域名
- POST /api/domains/:id/check - 检查域名证书
- PUT /api/domains/:id/auto-renewal - 切换自动续期状态

## 配置说明

### 后端配置
配置文件位置：`backend/configs/config.yaml`
```yaml
server:
  port: 8080
  host: "0.0.0.0"

mysql:
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_password"
  database: "ssl_monitor"

email:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  username: "your_email"
  password: "your_password"
  from_address: "your_email"
```

## 开发计划

- [x] 基础域名管理功能
- [x] 证书状态检查
- [ ] 邮件通知功能
- [ ] 自动续期功能
- [ ] 用户认证系统
- [ ] 操作日志记录
- [ ] 批量导入导出
- [ ] 证书详细信息显示

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支 (git checkout -b feature/AmazingFeature)
3. 提交您的更改 (git commit -m 'Add some AmazingFeature')
4. 推送到分支 (git push origin feature/AmazingFeature)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件