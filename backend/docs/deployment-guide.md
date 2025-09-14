# 部署指南

## 部署概述

本指南详细介绍了 WaitingToDo 项目的部署方法，包括开发环境搭建、生产环境部署和容器化部署等多种方式。

## 环境要求

### 基础环境

#### 后端环境
- **Go**: 1.24 或更高版本
- **MySQL**: 8.0 或更高版本
- **Redis**: 6.0 或更高版本
- **RabbitMQ**: 3.8 或更高版本
- **MinIO**: 最新稳定版本

#### 前端环境
- **Node.js**: 18.0 或更高版本
- **npm**: 9.0 或更高版本（或 yarn 1.22+）


## 开发环境部署

### 1. 环境准备

#### 安装 Go
```bash
# Linux/macOS
wget https://golang.org/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Windows
# 下载并安装 Go 安装包
# https://golang.org/dl/
```

#### 安装 Node.js
```bash
# 使用 nvm 安装（推荐）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18
nvm use 18

# 或直接下载安装包
# https://nodejs.org/
```

### 2. 数据库安装

#### MySQL 安装
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server
sudo mysql_secure_installation

# CentOS/RHEL
sudo yum install mysql-server
sudo systemctl start mysqld
sudo systemctl enable mysqld

# macOS
brew install mysql
brew services start mysql
```

#### Redis 安装
```bash
# Ubuntu/Debian
sudo apt install redis-server
sudo systemctl start redis
sudo systemctl enable redis

# CentOS/RHEL
sudo yum install redis
sudo systemctl start redis
sudo systemctl enable redis

# macOS
brew install redis
brew services start redis
```

#### RabbitMQ 安装
```bash
# Ubuntu/Debian
sudo apt install rabbitmq-server
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server

# 启用管理插件
sudo rabbitmq-plugins enable rabbitmq_management
```

#### MinIO 安装
```bash
# Linux
wget https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio
sudo mv minio /usr/local/bin/

# 启动 MinIO
minio server /data --console-address ":9001"
```

### 3. 项目部署

#### 克隆项目
```bash
git https://github.com/BlueOcean223/WaitingToDo.git
cd WaitingToDo
```

#### 后端部署
```bash
cd backend

# 安装依赖
go mod tidy

# 配置数据库
# 复制配置文件
cp config/config.yaml config/config.local.yaml

# 编辑配置文件
vim config/config.local.yaml
```

**配置文件示例：**
```yaml
mysql:
  dsn: "root:password@tcp(127.0.0.1:3306)/WaitingToDo?charset=utf8mb4&parseTime=True&loc=Local"
  max_idle_conns: 5
  max_open_conns: 20
  conn_max_lifetime: 1800
  conn_max_idle_time: 900

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

rabbitmq:
  dsn: "amqp://guest:guest@localhost:5672/"
  exchange: "social"
  queues:
    friend_request:
      name: "friend_requests"
      routing_key: "friend.request"
    team_request:
      name: "team_requests"
      routing_key: "team.request"

minio:
  endpoint: "127.0.0.1:9000"
  access_key_id: "minioadmin"
  secret_access_key: "minioadmin"

mail:
  smtp_host: "smtp.163.com"
  smtp_port: 465
  from: "your_email@163.com"
  password: "your_password"

log:
  level: "info"
  filename: "logs/app.log"
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true
  console: true
```

#### 数据库初始化
```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE WaitingToDo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;

# 运行项目（自动创建表）
go run cmd/main.go
```

#### 前端部署
```bash
cd frontend/vue

# 安装依赖
npm install
# 或使用 yarn
yarn install

# 配置环境变量
cp .env.development .env.local

# 编辑环境变量
vim .env.local
```

**环境变量示例：**
```bash
# .env.local
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=WaitingToDo
VITE_UPLOAD_URL=http://localhost:9000
```

#### 启动服务
```bash
# 启动后端（开发模式）
cd backend
go install github.com/cosmtrek/air@latest
air

# 启动前端（开发模式）
cd frontend/vue
npm run dev
```

## 生产环境部署

### 1. 服务器准备

#### 系统配置
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装必要工具
sudo apt install -y curl wget git vim htop

# 配置防火墙
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

#### 创建用户
```bash
# 创建应用用户
sudo useradd -m -s /bin/bash waitingtodo
sudo usermod -aG sudo waitingtodo

# 切换到应用用户
su - waitingtodo
```

### 2. 数据库部署

#### MySQL 生产配置
```bash
# 安装 MySQL
sudo apt install mysql-server

# 安全配置
sudo mysql_secure_installation

# 创建应用数据库和用户
mysql -u root -p
```

```sql
CREATE DATABASE WaitingToDo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'waitingtodo'@'localhost' IDENTIFIED BY 'strong_password';
GRANT ALL PRIVILEGES ON WaitingToDo.* TO 'waitingtodo'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

#### Redis 生产配置
```bash
# 安装 Redis
sudo apt install redis-server

# 配置 Redis
sudo vim /etc/redis/redis.conf
```

```conf
# /etc/redis/redis.conf
bind 127.0.0.1
port 6379
requirepass your_redis_password
maxmemory 256mb
maxmemory-policy allkeys-lru
```

```bash
# 重启 Redis
sudo systemctl restart redis
sudo systemctl enable redis
```

### 3. 应用部署

#### 后端部署
```bash
# 克隆代码
git clone https://github.com/your-username/WaitingToDo.git
cd WaitingToDo/backend

# 构建应用
go build -o waitingtodo cmd/main.go

# 创建配置文件
sudo mkdir -p /etc/waitingtodo
sudo cp config/config.yaml /etc/waitingtodo/
sudo chown waitingtodo:waitingtodo /etc/waitingtodo/config.yaml

# 创建日志目录
sudo mkdir -p /var/log/waitingtodo
sudo chown waitingtodo:waitingtodo /var/log/waitingtodo

# 创建 systemd 服务
sudo vim /etc/systemd/system/waitingtodo.service
```

**systemd 服务配置：**
```ini
[Unit]
Description=WaitingToDo Backend Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=waitingtodo
Group=waitingtodo
WorkingDirectory=/home/waitingtodo/WaitingToDo/backend
ExecStart=/home/waitingtodo/WaitingToDo/backend/waitingtodo
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

```bash
# 启动服务
sudo systemctl daemon-reload
sudo systemctl start waitingtodo
sudo systemctl enable waitingtodo

# 检查状态
sudo systemctl status waitingtodo
```

#### 前端部署
```bash
cd frontend/vue

# 安装依赖
npm install

# 构建生产版本
npm run build

# 部署到 Nginx
sudo cp -r dist/* /var/www/html/
sudo chown -R www-data:www-data /var/www/html/
```

### 4. Nginx 配置

#### 安装 Nginx
```bash
sudo apt install nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

#### 配置 Nginx
```bash
sudo vim /etc/nginx/sites-available/waitingtodo
```

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/html;
    index index.html;
    client_max_body_size 60M;

    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;
}
```

```bash
# 启用站点
sudo ln -s /etc/nginx/sites-available/waitingtodo /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 5. SSL 配置

#### 使用 Let's Encrypt
```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取 SSL 证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加以下行
0 12 * * * /usr/bin/certbot renew --quiet
```

## 容器化部署

### 1. Docker 部署

#### 安装 Docker
```bash
# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 创建 docker-compose.yml
```yaml
version: '3.8'

services:
  # MySQL 数据库
  mysql:
    image: mysql:8.0
    container_name: waitingtodo-mysql
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: WaitingToDo
      MYSQL_USER: waitingtodo
      MYSQL_PASSWORD: password
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    networks:
      - waitingtodo-network

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: waitingtodo-redis
    command: redis-server --requirepass password
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    networks:
      - waitingtodo-network

  # RabbitMQ 消息队列
  rabbitmq:
    image: rabbitmq:3-management
    container_name: waitingtodo-rabbitmq
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: password
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"
    networks:
      - waitingtodo-network

  # MinIO 对象存储
  minio:
    image: minio/minio:latest
    container_name: waitingtodo-minio
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    command: server /data --console-address ":9001"
    volumes:
      - minio_data:/data
    ports:
      - "9000:9000"
      - "9001:9001"
    networks:
      - waitingtodo-network

  # 后端服务
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: waitingtodo-backend
    depends_on:
      - mysql
      - redis
      - rabbitmq
      - minio
    environment:
      GIN_MODE: release
    volumes:
      - ./backend/config:/app/config
      - ./backend/logs:/app/logs
    ports:
      - "8080:8080"
    networks:
      - waitingtodo-network

  # 前端服务
  frontend:
    build:
      context: ./frontend/vue
      dockerfile: Dockerfile
    container_name: waitingtodo-frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - waitingtodo-network

volumes:
  mysql_data:
  redis_data:
  rabbitmq_data:
  minio_data:

networks:
  waitingtodo-network:
    driver: bridge
```

#### 启动容器
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend

# 停止服务
docker-compose down
```

### 2. Kubernetes 部署

#### 创建命名空间
```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: waitingtodo
```

#### 部署数据库
```yaml
# mysql-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
  namespace: waitingtodo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "rootpassword"
        - name: MYSQL_DATABASE
          value: "WaitingToDo"
        ports:
        - containerPort: 3306
        volumeMounts:
        - name: mysql-storage
          mountPath: /var/lib/mysql
      volumes:
      - name: mysql-storage
        persistentVolumeClaim:
          claimName: mysql-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
  namespace: waitingtodo
spec:
  selector:
    app: mysql
  ports:
  - port: 3306
    targetPort: 3306
```

#### 部署应用
```bash
# 应用所有配置
kubectl apply -f k8s/

# 查看部署状态
kubectl get pods -n waitingtodo
kubectl get services -n waitingtodo

# 查看日志
kubectl logs -f deployment/backend -n waitingtodo
```

## 监控和维护

### 1. 日志管理

#### 日志轮转
```bash
# 配置 logrotate
sudo vim /etc/logrotate.d/waitingtodo
```

```conf
/var/log/waitingtodo/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 waitingtodo waitingtodo
    postrotate
        systemctl reload waitingtodo
    endscript
}
```

### 2. 备份策略

#### 数据库备份
```bash
#!/bin/bash
# backup.sh
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/mysql"
mkdir -p $BACKUP_DIR

mysqldump -u waitingtodo -p WaitingToDo > $BACKUP_DIR/waitingtodo_$DATE.sql
gzip $BACKUP_DIR/waitingtodo_$DATE.sql

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
```

#### 定时备份
```bash
# 添加到 crontab
crontab -e
# 每天凌晨 2 点备份
0 2 * * * /home/waitingtodo/scripts/backup.sh
```

### 3. 性能监控

#### 系统监控
```bash
# 安装监控工具
sudo apt install htop iotop nethogs

# 查看系统资源
htop
iotop
nethogs

# 查看服务状态
sudo systemctl status waitingtodo
sudo journalctl -u waitingtodo -f
```

## 故障排除

### 常见问题

#### 1. 服务启动失败
```bash
# 查看服务状态
sudo systemctl status waitingtodo

# 查看详细日志
sudo journalctl -u waitingtodo -n 50

# 检查配置文件
sudo -u waitingtodo /home/waitingtodo/WaitingToDo/backend/waitingtodo --config-check
```

#### 2. 数据库连接问题
```bash
# 测试数据库连接
mysql -u waitingtodo -p -h localhost WaitingToDo

# 检查数据库状态
sudo systemctl status mysql

# 查看数据库日志
sudo tail -f /var/log/mysql/error.log
```

#### 3. 内存不足
```bash
# 查看内存使用
free -h

# 查看进程内存使用
ps aux --sort=-%mem | head

# 添加交换空间
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### 性能优化

#### 数据库优化
```sql
-- 查看慢查询
SHOW VARIABLES LIKE 'slow_query_log';
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;

-- 分析表
ANALYZE TABLE tasks;
ANALYZE TABLE users;

-- 优化表
OPTIMIZE TABLE tasks;
```

#### 应用优化
```bash
# 调整 Go 运行时参数
export GOGC=100
export GOMAXPROCS=4

# 调整系统参数
echo 'net.core.somaxconn = 1024' | sudo tee -a /etc/sysctl.conf
echo 'net.ipv4.tcp_max_syn_backlog = 1024' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

---

*部署完成后，请确保定期更新系统和应用，监控系统性能，及时处理安全更新。*