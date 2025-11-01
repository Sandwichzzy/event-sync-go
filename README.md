好的！我已经检查了你的系统，PostgreSQL 确实还未安装。我为你准备了完整的安装和迁移步骤。

PostgreSQL 安装和数据迁移步骤

根据你的项目配置（.env 文件），项目已经配置好了 PostgreSQL 连接参数：
- 数据库名: event_sync
- 用户名: Sandwich
- 密码: zzy15198
- 端口: 5432

步骤 1: 安装 PostgreSQL

由于你使用的是 Ubuntu/Debian 系统，我们可以使用 apt 安装：

# 更新软件包列表
sudo apt update

# 安装 PostgreSQL 和相关工具
sudo apt install postgresql postgresql-contrib -y

你想现在开始安装吗？我可以帮你执行这些命令。

步骤 2: 启动并配置 PostgreSQL 服务

# 启动 PostgreSQL 服务
sudo systemctl start postgresql

# 设置开机自启
sudo systemctl enable postgresql

# 检查服务状态
sudo systemctl status postgresql

步骤 3: 创建数据库用户和数据库

# 切换到 postgres 用户并创建新用户
sudo -u postgres psql -c "CREATE USER \"Sandwich\" WITH PASSWORD 'zzy15198';"

# 创建数据库并设置所有者
sudo -u postgres psql -c "CREATE DATABASE event_sync OWNER \"Sandwich\";"

# 授予用户所有权限
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE event_sync TO \"Sandwich\";"

步骤 4: 运行数据库迁移脚本

你的项目已经有迁移文件 migrations/0001_create_schema.sql，它会创建以下表：
- block_headers - 区块头信息
- contract_events - 合约事件
- deposit_tokens - 存款代币
- withdraw_tokens - 提款代币
- grant_reward_tokens - 授予奖励代币
- withdraw_manager_update - 提款管理器更新

运行迁移（两种方式）：

方式 A: 使用项目的迁移功能
# 先加载环境变量
source .env

# 然后运行你的应用程序，它会自动执行迁移
./event-sync migrate  # 或者使用你的应用程序的迁移命令

方式 B: 手动运行 SQL
# 直接执行 SQL 文件
PGPASSWORD='zzy15198' psql -h 127.0.0.1 -U Sandwich -d event_sync -f ./migrations/0001_create_schema.sql

步骤 5: 验证安装

# 连接到数据库并查看表
PGPASSWORD='zzy15198' psql -h 127.0.0.1 -U Sandwich -d event_sync -c "\dt"

你想让我帮你执行这些步骤吗？我可以逐步为你运行这些命令。
