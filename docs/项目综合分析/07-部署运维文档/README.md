# 7. 部署运维文档

## 7.1 部署概述

### 7.1.1 部署架构

葛洲坝船闸导航系统支持多种部署模式：

```text
┌─────────────────────────────────────────────────────────┐
│                    生产环境                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 主服务器     │ │ 备用服务器   │ │ 监控服务器   │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    网络层                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 负载均衡器   │ │ 防火墙       │ │ 交换机       │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    设备层                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 雷达设备     │ │ 云台设备     │ │ LED设备      │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
```

### 7.1.2 系统要求

**硬件要求：**

- CPU: 4核心以上
- 内存: 8GB以上
- 存储: 100GB以上SSD
- 网络: 千兆以太网

**软件要求：**

- 操作系统: Ubuntu 20.04 LTS / CentOS 8
- Go版本: 1.24.3+
- MySQL: 8.0+
- NATS: 2.11.4+

## 7.2 部署指南

### 7.2.1 环境准备

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装Go
wget https://golang.org/dl/go1.24.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装MySQL
sudo apt install mysql-server -y
sudo systemctl start mysql
sudo systemctl enable mysql

# 安装NATS
wget https://github.com/nats-io/nats-server/releases/download/v2.11.4/nats-server-v2.11.4-linux-amd64.tar.gz
tar -xzf nats-server-v2.11.4-linux-amd64.tar.gz
sudo mv nats-server-v2.11.4-linux-amd64/nats-server /usr/local/bin/
```

### 7.2.2 应用部署

```bash
# 克隆代码
git clone https://github.com/your-org/navigate.git
cd navigate

# 构建应用
go build -o navlock-server cmd/server/navlock-gezhouba/main.go

# 创建配置目录
sudo mkdir -p /etc/navlock
sudo mkdir -p /var/lib/navlock
sudo mkdir -p /var/log/navlock

# 复制配置文件
sudo cp config/config.prod.yaml /etc/navlock/config.yaml

# 创建系统服务
sudo tee /etc/systemd/system/navlock.service > /dev/null <<EOF
[Unit]
Description=Navlock Service
After=network.target mysql.service

[Service]
Type=simple
User=navlock
Group=navlock
WorkingDirectory=/opt/navlock
ExecStart=/opt/navlock/navlock-server
Restart=always
RestartSec=5
Environment=NAVLOCK_ENV=production

[Install]
WantedBy=multi-user.target
EOF

# 创建用户
sudo useradd -r -s /bin/false navlock
sudo chown -R navlock:navlock /var/lib/navlock /var/log/navlock

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable navlock
sudo systemctl start navlock
```

### 7.2.3 数据库初始化

```sql
-- 创建数据库
CREATE DATABASE navlock_prod CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE navlock_web CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户
CREATE USER 'navlock'@'localhost' IDENTIFIED BY 'secure_password';
CREATE USER 'navlock_web'@'localhost' IDENTIFIED BY 'secure_password';

-- 授权
GRANT ALL PRIVILEGES ON navlock_prod.* TO 'navlock'@'localhost';
GRANT ALL PRIVILEGES ON navlock_web.* TO 'navlock_web'@'localhost';
FLUSH PRIVILEGES;

-- 导入初始数据
mysql -u navlock -p navlock_prod < docs/gezhouba/系统架构图/Dump20230301.sql
```

## 7.3 运维手册

### 7.3.1 服务管理

```bash
# 查看服务状态
sudo systemctl status navlock

# 启动服务
sudo systemctl start navlock

# 停止服务
sudo systemctl stop navlock

# 重启服务
sudo systemctl restart navlock

# 查看日志
sudo journalctl -u navlock -f

# 查看实时日志
sudo tail -f /var/log/navlock/navlock.log
```

### 7.3.2 监控检查

```bash
# 检查系统资源
htop
df -h
free -h

# 检查网络连接
netstat -tulpn | grep navlock
ss -tulpn | grep navlock

# 检查数据库连接
mysql -u navlock -p -e "SHOW PROCESSLIST;"

# 检查NATS连接
nats-server --version
```

### 7.3.3 备份恢复

```bash
# 数据库备份
mysqldump -u navlock -p navlock_prod > backup_$(date +%Y%m%d_%H%M%S).sql
mysqldump -u navlock_web -p navlock_web > backup_web_$(date +%Y%m%d_%H%M%S).sql

# 配置文件备份
sudo cp /etc/navlock/config.yaml /etc/navlock/config.yaml.backup

# 应用备份
sudo cp /opt/navlock/navlock-server /opt/navlock/navlock-server.backup

# 数据恢复
mysql -u navlock -p navlock_prod < backup_20240101_120000.sql
```

## 7.4 故障处理

### 7.4.1 常见故障

**服务无法启动：**

```bash
# 检查配置文件
sudo navlock-server --config-check

# 检查端口占用
sudo netstat -tulpn | grep :8080

# 检查权限
ls -la /var/lib/navlock /var/log/navlock
```

**数据库连接失败：**

```bash
# 检查MySQL服务
sudo systemctl status mysql

# 检查连接配置
mysql -u navlock -p -h localhost

# 检查防火墙
sudo ufw status
```

**设备通信异常：**

```bash
# 检查网络连接
ping 192.168.1.100

# 检查设备状态
telnet 192.168.1.100 8001

# 查看设备日志
sudo journalctl -u navlock | grep radar
```

### 7.4.2 故障诊断

```bash
# 系统诊断脚本
#!/bin/bash
echo "=== 系统诊断报告 ==="
echo "时间: $(date)"
echo ""

echo "=== 系统资源 ==="
echo "CPU使用率:"
top -bn1 | grep "Cpu(s)" | awk '{print $2}'
echo "内存使用:"
free -h
echo "磁盘使用:"
df -h

echo ""
echo "=== 服务状态 ==="
sudo systemctl status navlock --no-pager

echo ""
echo "=== 网络连接 ==="
netstat -tulpn | grep navlock

echo ""
echo "=== 最近日志 ==="
sudo journalctl -u navlock --no-pager -n 50
```

## 7.5 监控告警

### 7.5.1 监控指标

**系统指标：**

- CPU使用率
- 内存使用率
- 磁盘使用率
- 网络流量

**应用指标：**

- 服务响应时间
- 错误率
- 并发连接数
- 消息处理量

**业务指标：**

- 船舶通过数量
- 设备在线率
- 异常事件数量
- 系统可用性

### 7.5.2 告警配置

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'navlock'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s

# alertmanager.yml
route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'web.hook'

receivers:
  - name: 'web.hook'
    webhook_configs:
      - url: 'http://127.0.0.1:5001/'
```

### 7.5.3 告警规则

```yaml
# rules.yml
groups:
  - name: navlock.rules
    rules:
      - alert: HighCPUUsage
        expr: cpu_usage > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "CPU使用率过高"
          description: "CPU使用率超过80%持续5分钟"

      - alert: HighMemoryUsage
        expr: memory_usage > 85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "内存使用率过高"
          description: "内存使用率超过85%持续5分钟"

      - alert: ServiceDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "服务不可用"
          description: "Navlock服务已停止运行"
```

## 7.6 性能优化

### 7.6.1 系统优化

```bash
# 内核参数优化
echo 'net.core.somaxconn = 65535' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_max_syn_backlog = 65535' >> /etc/sysctl.conf
echo 'vm.swappiness = 10' >> /etc/sysctl.conf
sysctl -p

# 文件描述符限制
echo '* soft nofile 65535' >> /etc/security/limits.conf
echo '* hard nofile 65535' >> /etc/security/limits.conf
```

### 7.6.2 数据库优化

```sql
-- MySQL配置优化
SET GLOBAL innodb_buffer_pool_size = 1073741824; -- 1GB
SET GLOBAL innodb_log_file_size = 268435456; -- 256MB
SET GLOBAL innodb_flush_log_at_trx_commit = 2;
SET GLOBAL innodb_flush_method = O_DIRECT;

-- 创建索引
CREATE INDEX idx_ship_timestamp ON ship_positions(timestamp);
CREATE INDEX idx_event_type ON events(event_type);
```

### 7.6.3 应用优化

```go
// 连接池配置
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)

// 缓存配置
cache := cache.New(5*time.Minute, 10*time.Minute)
```

## 7.7 安全配置

### 7.7.1 网络安全

```bash
# 防火墙配置
sudo ufw enable
sudo ufw allow ssh
sudo ufw allow 8080/tcp
sudo ufw allow 4222/tcp
sudo ufw allow 3306/tcp

# SSL/TLS配置
sudo apt install certbot -y
sudo certbot certonly --standalone -d your-domain.com
```

### 7.7.2 访问控制

```bash
# SSH安全配置
sudo sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sudo sed -i 's/PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo systemctl restart sshd

# 文件权限
sudo chmod 600 /etc/navlock/config.yaml
sudo chown navlock:navlock /etc/navlock/config.yaml
```

### 7.7.3 数据安全

```bash
# 数据加密
sudo apt install cryptsetup -y
sudo cryptsetup luksFormat /dev/sdb1
sudo cryptsetup luksOpen /dev/sdb1 encrypted_data

# 备份加密
gpg --symmetric --cipher-algo AES256 backup.sql
```

## 7.8 升级维护

### 7.8.1 版本升级

```bash
# 备份当前版本
sudo systemctl stop navlock
sudo cp /opt/navlock/navlock-server /opt/navlock/navlock-server.old

# 部署新版本
sudo cp navlock-server.new /opt/navlock/navlock-server
sudo chown navlock:navlock /opt/navlock/navlock-server
sudo chmod +x /opt/navlock/navlock-server

# 启动新版本
sudo systemctl start navlock
sudo systemctl status navlock

# 验证功能
curl http://localhost:8080/health
```

### 7.8.2 回滚操作

```bash
# 停止服务
sudo systemctl stop navlock

# 恢复旧版本
sudo cp /opt/navlock/navlock-server.old /opt/navlock/navlock-server
sudo chown navlock:navlock /opt/navlock/navlock-server
sudo chmod +x /opt/navlock/navlock-server

# 启动服务
sudo systemctl start navlock
```

### 7.8.3 维护计划

**日常维护：**

- 检查系统资源使用情况
- 查看服务运行状态
- 备份重要数据
- 清理日志文件

**周维护：**

- 更新系统安全补丁
- 检查数据库性能
- 验证备份完整性
- 分析系统日志

**月维护：**

- 全面系统检查
- 性能优化调整
- 安全审计
- 容量规划
